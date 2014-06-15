package dl

// #include <stdlib.h>
// #include "trampoline.h"
import "C"

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

type rFunc func([]reflect.Value) []reflect.Value

func makeTrampoline(typ reflect.Type, handle unsafe.Pointer) (rFunc, error) {
	numOut := typ.NumOut()
	if numOut > 1 {
		return nil, fmt.Errorf("C functions can return 0 or 1 values, not %d", numOut)
	}
	var out reflect.Type
	var kind reflect.Kind
	outFlag := C.int(0)
	if numOut == 1 {
		out = typ.Out(0)
		kind = out.Kind()
		if kind == reflect.Float32 || kind == reflect.Float64 {
			outFlag |= C.ARG_FLAG_FLOAT
		}
	}
	return func(in []reflect.Value) []reflect.Value {
		count := len(in)
		args := make([]unsafe.Pointer, count)
		flags := make([]C.int, count+1)
		flags[count] = outFlag
		for ii, v := range in {
			switch v.Kind() {
			case reflect.String:
				s := C.CString(v.String())
				defer C.free(unsafe.Pointer(s))
				args[ii] = unsafe.Pointer(s)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				args[ii] = unsafe.Pointer(uintptr(v.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
			case reflect.Float32:
				args[ii] = unsafe.Pointer(uintptr(math.Float32bits(float32(v.Float()))))
				flags[ii] |= C.ARG_FLAG_FLOAT
			case reflect.Float64:
				args[ii] = unsafe.Pointer(uintptr(math.Float64bits(v.Float())))
				flags[ii] |= C.ARG_FLAG_FLOAT
			case reflect.Slice:
				if v.Len() > 0 {
					args[ii] = unsafe.Pointer(v.Index(0).UnsafeAddr())
				}
			default:
				panic(fmt.Errorf("can't bind value of type %s", v.Type()))
			}
		}
		ret := C.call(handle, &args[0], &flags[0], C.int(count))
		if numOut > 0 {
			var v reflect.Value
			switch kind {
			case reflect.Int:
				v = reflect.ValueOf(int(int32(uintptr(ret))))
			case reflect.Float32:
				v = reflect.ValueOf(math.Float32frombits(uint32(uintptr(ret))))
			case reflect.Float64:
				v = reflect.ValueOf(math.Float64frombits(uint64(uintptr(ret))))
			default:
				panic(fmt.Errorf("can't retrieve value of type %s", out))
			}
			return []reflect.Value{v}
		}
		return nil
	}, nil
}
