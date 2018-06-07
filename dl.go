package dl

/*
#cgo LDFLAGS: -ldl

#define _GNU_SOURCE
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// RTLD_* values. See man dlopen.
const (
	RTLD_LAZY     = int(C.RTLD_LAZY)
	RTLD_NOW      = int(C.RTLD_NOW)
	RTLD_GLOBAL   = int(C.RTLD_GLOBAL)
	RTLD_LOCAL    = int(C.RTLD_LOCAL)
	RTLD_NODELETE = int(C.RTLD_NODELETE)
	RTLD_NOLOAD   = int(C.RTLD_NOLOAD)
)

// mu is the dl* call mutex.
var mu sync.Mutex

// dlerror wraps C.dlerror, returning a error from C.dlerror.
func dlerror() error {
	return errors.New(C.GoString(C.dlerror()))
}

// dlopen wraps C.dlopen, opening a handle for library n, passing flags.
func dlopen(name string, flags int) (unsafe.Pointer, error) {
	mu.Lock()
	defer mu.Unlock()

	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	h := C.dlopen(n, C.int(flags))
	if h == nil {
		return nil, dlerror()
	}
	return h, nil
}

// dlclose wraps C.dlclose, closing handle l.
func dlclose(h unsafe.Pointer) error {
	mu.Lock()
	defer mu.Unlock()

	if C.dlclose(h) != 0 {
		return dlerror()
	}
	return nil
}

// dlsym wraps C.dlsym, loading from handle h the symbol name n, and returning a
// pointer to the loaded symbol.
func dlsym(h unsafe.Pointer, n *C.char) (unsafe.Pointer, error) {
	mu.Lock()
	defer mu.Unlock()

	sym := C.dlsym(h, n)
	if sym == nil {
		return nil, dlerror()
	}
	return sym, nil
}

// cast converts sym into v.
func cast(sym unsafe.Pointer, v interface{}) error {
	val := reflect.ValueOf(v)
	if !val.IsValid() || val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("v must be a pointer and cannot be nil")
	}

	elem := val.Elem()
	switch elem.Kind() {
	// treat Go int/uint as long, since it depends on the platform
	case reflect.Int:
		elem.SetInt(int64(*(*int)(sym)))
	case reflect.Uint:
		elem.SetUint(uint64(*(*uint)(sym)))

	case reflect.Int8:
		elem.SetInt(int64(*(*int8)(sym)))
	case reflect.Int16:
		elem.SetInt(int64(*(*int16)(sym)))
	case reflect.Int32:
		elem.SetInt(int64(*(*int32)(sym)))
	case reflect.Int64:
		elem.SetInt(int64(*(*int64)(sym)))

	case reflect.Uint8:
		elem.SetUint(uint64(*(*uint8)(sym)))
	case reflect.Uint16:
		elem.SetUint(uint64(*(*uint16)(sym)))
	case reflect.Uint32:
		elem.SetUint(uint64(*(*uint32)(sym)))
	case reflect.Uint64:
		elem.SetUint(uint64(*(*uint64)(sym)))

	case reflect.Float32:
		elem.SetFloat(float64(*(*float32)(sym)))
	case reflect.Float64:
		elem.SetFloat(float64(*(*float64)(sym)))

	case reflect.Uintptr:
		elem.SetUint(uint64(*(*uintptr)(sym)))
	case reflect.Ptr:
		v := reflect.NewAt(elem.Type().Elem(), sym)
		elem.Set(v)
	case reflect.UnsafePointer:
		elem.SetPointer(sym)

	case reflect.String:
		elem.SetString(C.GoString(*(**C.char)(sym)))

	case reflect.Func:
		typ := elem.Type()
		tr, err := makeTrampoline(typ, sym)
		if err != nil {
			return err
		}
		v := reflect.MakeFunc(typ, tr)
		elem.Set(v)

	default:
		return fmt.Errorf("cannot convert to type %T", elem.Kind())
	}

	return nil
}

// Sym wraps loading symbol name from handle h, decoding the value to v.
func Sym(h unsafe.Pointer, name string, v interface{}) error {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))

	sym, err := dlsym(h, n)
	if err != nil {
		return err
	}

	return cast(sym, v)
}

// SymDefault loads symbol name into v from the RTLD_DEFAULT handle.
func SymDefault(name string, v interface{}) error {
	return Sym(C.RTLD_DEFAULT, name, v)
}

// SymNext loads symbol name into v from the RTLD_NEXT handle.
func SymNext(name string, v interface{}) error {
	return Sym(C.RTLD_NEXT, name, v)
}
