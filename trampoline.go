package dl

// #include <stdlib.h>
// #include "trampoline.h"
import "C"

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

var (
	emptyType = reflect.TypeOf((*interface{})(nil)).Elem()
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
		if typ.IsVariadic() && len(in) > 0 {
			last := in[len(in)-1]
			in = in[:len(in)-1]
			if last.Len() > 0 {
				for ii := 0; ii < last.Len(); ii++ {
					in = append(in, last.Index(ii))
				}
			}
		}
		count := len(in)
		args := make([]unsafe.Pointer, count)
		flags := make([]C.int, count+1)
		flags[count] = outFlag
		for ii, v := range in {
			if v.Type() == emptyType {
				v = reflect.ValueOf(v.Interface())
			}
			switch v.Kind() {
			case reflect.String:
				s := C.CString(v.String())
				defer C.free(unsafe.Pointer(s))
				args[ii] = unsafe.Pointer(s)
				flags[ii] |= C.ARG_FLAG_SIZE_PTR
			case reflect.Int:
				args[ii] = unsafe.Pointer(uintptr(v.Int()))
				if v.Type().Size() == 4 {
					flags[ii] = C.ARG_FLAG_SIZE_32
				} else {
					flags[ii] = C.ARG_FLAG_SIZE_64
				}
			case reflect.Int8:
				args[ii] = unsafe.Pointer(uintptr(v.Int()))
				flags[ii] = C.ARG_FLAG_SIZE_8
			case reflect.Int16:
				args[ii] = unsafe.Pointer(uintptr(v.Int()))
				flags[ii] = C.ARG_FLAG_SIZE_16
			case reflect.Int32:
				args[ii] = unsafe.Pointer(uintptr(v.Int()))
				flags[ii] = C.ARG_FLAG_SIZE_32
			case reflect.Int64:
				args[ii] = unsafe.Pointer(uintptr(v.Int()))
				flags[ii] = C.ARG_FLAG_SIZE_64
			case reflect.Uint:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
				if v.Type().Size() == 4 {
					flags[ii] = C.ARG_FLAG_SIZE_32
				} else {
					flags[ii] = C.ARG_FLAG_SIZE_64
				}
			case reflect.Uint8:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
				flags[ii] = C.ARG_FLAG_SIZE_8
			case reflect.Uint16:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
				flags[ii] = C.ARG_FLAG_SIZE_16
			case reflect.Uint32:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
				flags[ii] = C.ARG_FLAG_SIZE_32
			case reflect.Uint64:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
				flags[ii] = C.ARG_FLAG_SIZE_64
			case reflect.Float32:
				args[ii] = unsafe.Pointer(uintptr(math.Float32bits(float32(v.Float()))))
				flags[ii] |= C.ARG_FLAG_FLOAT | C.ARG_FLAG_SIZE_32
			case reflect.Float64:
				args[ii] = unsafe.Pointer(uintptr(math.Float64bits(v.Float())))
				flags[ii] |= C.ARG_FLAG_FLOAT | C.ARG_FLAG_SIZE_64
			case reflect.Ptr:
				args[ii] = unsafe.Pointer(v.Pointer())
				flags[ii] |= C.ARG_FLAG_SIZE_PTR
			case reflect.Slice:
				if v.Len() > 0 {
					args[ii] = unsafe.Pointer(v.Index(0).UnsafeAddr())
				}
				flags[ii] |= C.ARG_FLAG_SIZE_PTR
			case reflect.Uintptr:
				args[ii] = unsafe.Pointer(uintptr(v.Uint()))
				flags[ii] |= C.ARG_FLAG_SIZE_PTR
			default:
				panic(fmt.Errorf("can't bind value of type %s", v.Type()))
			}
		}
		var argp *unsafe.Pointer
		if count > 0 {
			argp = &args[0]
		}
		var ret unsafe.Pointer
		if C.call(handle, argp, &flags[0], C.int(count), &ret) != 0 {
			s := C.GoString((*C.char)(ret))
			C.free(ret)
			panic(errors.New(s))
		}
		if numOut > 0 {
			var v reflect.Value
			switch kind {
			case reflect.Int:
				v = reflect.ValueOf(int(uintptr(ret)))
			case reflect.Int8:
				v = reflect.ValueOf(int8(uintptr(ret)))
			case reflect.Int16:
				v = reflect.ValueOf(int16(uintptr(ret)))
			case reflect.Int32:
				v = reflect.ValueOf(int32(uintptr(ret)))
			case reflect.Int64:
				v = reflect.ValueOf(int64(uintptr(ret)))
			case reflect.Uint:
				v = reflect.ValueOf(uint(uintptr(ret)))
			case reflect.Uint8:
				v = reflect.ValueOf(uint8(uintptr(ret)))
			case reflect.Uint16:
				v = reflect.ValueOf(uint16(uintptr(ret)))
			case reflect.Uint32:
				v = reflect.ValueOf(uint32(uintptr(ret)))
			case reflect.Uint64:
				v = reflect.ValueOf(uint64(uintptr(ret)))
			case reflect.Float32:
				v = reflect.ValueOf(math.Float32frombits(uint32(uintptr(ret))))
			case reflect.Float64:
				v = reflect.ValueOf(math.Float64frombits(uint64(uintptr(ret))))
			case reflect.Ptr:
				if out.Elem().Kind() == reflect.String && ret != nil {
					s := C.GoString((*C.char)(ret))
					v = reflect.ValueOf(&s)
					break
				}
				v = reflect.NewAt(out.Elem(), ret)
			case reflect.String:
				s := C.GoString((*C.char)(ret))
				v = reflect.ValueOf(s)
			case reflect.Uintptr:
				v = reflect.ValueOf(uintptr(ret))
			case reflect.UnsafePointer:
				v = reflect.ValueOf(ret)
			default:
				panic(fmt.Errorf("can't retrieve value of type %s", out))
			}
			return []reflect.Value{v}
		}
		return nil
	}, nil
}
