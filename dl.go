package dlsym

// #include <dlfcn.h>
// #include <stdlib.h>
// #cgo LDFLAGS: -ldl
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

const (
	RTLD_LAZY = int(C.RTLD_LAZY)
	RTLD_NOW  = int(C.RTLD_NOW)
)

var (
	mu sync.Mutex
)

type DL struct {
	mu     sync.Mutex
	handle unsafe.Pointer
}

func Open(name string, flag int) (*DL, error) {
	if flag&RTLD_LAZY == 0 && flag&RTLD_NOW == 0 {
		flag |= RTLD_NOW
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	mu.Lock()
	defer mu.Unlock()
	handle := C.dlopen(s, C.int(flag))
	if handle == nil {
		return nil, dlerror()
	}
	return &DL{
		handle: handle,
	}, nil
}

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
		// C int is 32 bits even in 64 bit ABI
		elem.SetInt(int64(*(*int32)(handle)))
	case reflect.Int8:
		elem.SetInt(int64(*(*int8)(handle)))
	case reflect.Int16:
		elem.SetInt(int64(*(*int16)(handle)))
	case reflect.Int32:
		elem.SetInt(int64(*(*int32)(handle)))
	case reflect.Int64:
		elem.SetInt(int64(*(*int64)(handle)))
	case reflect.Uint:
		// C uint is 32 bits even in 64 bit ABI
		elem.SetUint(uint64(*(*uint32)(handle)))
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
	case reflect.String:
		elem.SetString(C.GoString(*(**C.char)(handle)))
	case reflect.UnsafePointer:
		elem.SetPointer(handle)
	default:
		return fmt.Errorf("invalid out type %T", out)
	}
	return nil
}

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
