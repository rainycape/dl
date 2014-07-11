package dl

// #include <dlfcn.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -ldl
import "C"

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

const (
	// dlopen() flags. See man dlopen.
	RTLD_LAZY     = int(C.RTLD_LAZY)
	RTLD_NOW      = int(C.RTLD_NOW)
	RTLD_GLOBAL   = int(C.RTLD_GLOBAL)
	RTLD_LOCAL    = int(C.RTLD_LOCAL)
	RTLD_NODELETE = int(C.RTLD_NODELETE)
	RTLD_NOLOAD   = int(C.RTLD_NOLOAD)
)

var (
	mu sync.Mutex
)

// DL represents an opened dynamic library. Use Open
// to initialize a DL and use DL.Close when you're finished
// with it. Note that when the DL is closed all its loaded
// symbols become invalid.
type DL struct {
	mu     sync.Mutex
	handle unsafe.Pointer
}

// Open opens the shared library identified by the given name
// with the given flags. See man dlopen for the available flags
// and its meaning. Note that the only difference with dlopen is that
// if nor RTLD_LAZY nor RTLD_NOW are specified, Open defaults to
// RTLD_NOW rather than returning an error. If the name argument
// passed to name does not have extension, the default for the
// platform will be appended to it (e.g. .so, .dylib, etc...).
func Open(name string, flag int) (*DL, error) {
	if flag&RTLD_LAZY == 0 && flag&RTLD_NOW == 0 {
		flag |= RTLD_NOW
	}
	if name != "" && filepath.Ext(name) == "" {
		name = name + LibExt
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	mu.Lock()
	handle := C.dlopen(s, C.int(flag))
	var err error
	if handle == nil {
		err = dlerror()
	}
	mu.Unlock()
	if err != nil {
		if runtime.GOOS == "linux" && name == "libc.so" {
			// In most distros libc.so is now a text file
			// and in order to dlopen() it the name libc.so.6
			// must be used.
			return Open(name+".6", flag)
		}
		return nil, err
	}
	return &DL{
		handle: handle,
	}, nil
}

// Sym loads the symbol identified by the given name into
// the out parameter. Note that out must always be a pointer.
// See the package documentation to learn how types are mapped
// between Go and C.
func (d *DL) Sym(symbol string, out interface{}) error {
	s := C.CString(symbol)
	defer C.free(unsafe.Pointer(s))
	mu.Lock()
	handle := C.dlsym(d.handle, s)
	if handle == nil {
		err := dlerror()
		mu.Unlock()
		return err
	}
	mu.Unlock()
	val := reflect.ValueOf(out)
	if !val.IsValid() || val.Kind() != reflect.Ptr {
		return fmt.Errorf("out must be a pointer, not %T", out)
	}
	if val.IsNil() {
		return errors.New("out can't be nil")
	}
	elem := val.Elem()
	switch elem.Kind() {
	case reflect.Int:
		// We treat Go's int as long, since it
		// varies depending on the platform bit size
		elem.SetInt(int64(*(*int)(handle)))
	case reflect.Int8:
		elem.SetInt(int64(*(*int8)(handle)))
	case reflect.Int16:
		elem.SetInt(int64(*(*int16)(handle)))
	case reflect.Int32:
		elem.SetInt(int64(*(*int32)(handle)))
	case reflect.Int64:
		elem.SetInt(int64(*(*int64)(handle)))
	case reflect.Uint:
		// We treat Go's uint as unsigned long, since it
		// varies depending on the platform bit size
		elem.SetUint(uint64(*(*uint)(handle)))
	case reflect.Uint8:
		elem.SetUint(uint64(*(*uint8)(handle)))
	case reflect.Uint16:
		elem.SetUint(uint64(*(*uint16)(handle)))
	case reflect.Uint32:
		elem.SetUint(uint64(*(*uint32)(handle)))
	case reflect.Uint64:
		elem.SetUint(uint64(*(*uint64)(handle)))
	case reflect.Uintptr:
		elem.SetUint(uint64(*(*uintptr)(handle)))
	case reflect.Float32:
		elem.SetFloat(float64(*(*float32)(handle)))
	case reflect.Float64:
		elem.SetFloat(float64(*(*float64)(handle)))
	case reflect.Func:
		typ := elem.Type()
		tr, err := makeTrampoline(typ, handle)
		if err != nil {
			return err
		}
		v := reflect.MakeFunc(typ, tr)
		elem.Set(v)
	case reflect.Ptr:
		v := reflect.NewAt(elem.Type().Elem(), handle)
		elem.Set(v)
	case reflect.String:
		elem.SetString(C.GoString(*(**C.char)(handle)))
	case reflect.UnsafePointer:
		elem.SetPointer(handle)
	default:
		return fmt.Errorf("invalid out type %T", out)
	}
	return nil
}

// Close closes the shared library handle. All symbols
// loaded from the library will become invalid.
func (d *DL) Close() error {
	if d.handle != nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if d.handle != nil {
			mu.Lock()
			defer mu.Unlock()
			if C.dlclose(d.handle) != 0 {
				return dlerror()
			}
			d.handle = nil
		}
	}
	return nil
}

func dlerror() error {
	s := C.dlerror()
	return errors.New(C.GoString(s))
}
