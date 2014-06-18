package dl

import (
	"os/exec"
	"path/filepath"
	"testing"
	"unsafe"
)

var (
	testLib = filepath.Join("testdata", "lib")
)

func openTestLib(t *testing.T) *DL {
	dl, err := Open(testLib, 0)
	if err != nil {
		t.Fatal(err)
	}
	if testing.Verbose() {
		var verbose *int32
		if err := dl.Sym("verbose", &verbose); err != nil {
			t.Fatal(err)
		}
		*verbose = 1
	}
	return dl
}

func TestOpen(t *testing.T) {
	dl, err := Open("libc", 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := dl.Close(); err != nil {
		t.Fatal(err)
	}
	// Test double close
	if err := dl.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestDlsymVars(t *testing.T) {
	dl := openTestLib(t)
	defer dl.Close()
	var (
		s    string
		c    byte
		u32  uint32
		i    int32
		lng  int
		ptr  uintptr
		uptr unsafe.Pointer
	)
	if err := dl.Sym("my_string", &s); err != nil {
		t.Error(err)
	} else if s != "mystring" {
		t.Errorf("expecting s = \"mystring\", got %q instead", s)
	}
	if err := dl.Sym("my_char", &c); err != nil {
		t.Error(err)
	} else if c != 42 {
		t.Errorf("expecting ch = 42, got %v instead", c)
	}
	if err := dl.Sym("my_uint32", &u32); err != nil {
		t.Error(err)
	} else if u32 != 1337 {
		t.Errorf("expecting u32 = 1337, got %v instead", u32)
	}
	if err := dl.Sym("my_int", &i); err != nil {
		t.Error(err)
	} else if i != -9000 {
		t.Errorf("expecting i = -9000, got %v instead", i)
	}
	if err := dl.Sym("my_long", &lng); err != nil {
		t.Error(err)
	} else if lng != -9000 {
		t.Errorf("expecting lng = -9000, got %v instead", i)
	}
	if err := dl.Sym("my_pointer", &ptr); err != nil {
		t.Error(err)
	} else if ptr != 0xdeadbeef {
		t.Errorf("expecting ptr = 0xdeadbeef, got 0x%x instead", ptr)
	}
	// unsafe.Pointer return the address of the symbol
	if err := dl.Sym("my_pointer", &uptr); err != nil {
		t.Error(err)
	} else if uintptr(uptr) == 0xdeadbeef {
		t.Errorf("expecting uptr != 0xdeadbeef, got 0x%x instead", uintptr(uptr))
	}
}

func TestCounter(t *testing.T) {
	dl := openTestLib(t)
	defer dl.Close()
	// C's int is always 32 bits
	var counterPointer *int32
	if err := dl.Sym("counter", &counterPointer); err != nil {
		t.Fatal(err)
	}
	var increaseCounter func()
	if err := dl.Sym("increase_counter", &increaseCounter); err != nil {
		t.Fatal(err)
	}
	increaseCounter()
	var counter int32
	if err := dl.Sym("counter", &counter); err != nil {
		t.Fatal(err)
	}
	if counter != 1 || *counterPointer != 1 {
		t.Fatalf("counter should be 1, it's (%d / %d)", counter, *counterPointer)
	}
	increaseCounter()
	if err := dl.Sym("counter", &counter); err != nil {
		t.Fatal(err)
	}
	if counter != 2 || *counterPointer != 2 {
		t.Fatalf("counter should be 2, it's (%d / %d)", counter, *counterPointer)
	}
}

func TestStrlen(t *testing.T) {
	dl, err := Open("libc", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer dl.Close()
	const s = "golang"
	var strlen func(string) int32
	if err := dl.Sym("strlen", &strlen); err != nil {
		t.Error(err)
	} else if int(strlen(s)) != len(s) {
		t.Errorf("expecting strlen(%q) = %v, got %v instead", s, len(s), strlen(s))
	}
}

func TestFunctions(t *testing.T) {
	dl := openTestLib(t)
	defer dl.Close()
	var square func(float64) float64
	var squaref func(float32) float32

	if err := dl.Sym("square", &square); err != nil {
		t.Fatal(err)
	}

	if err := dl.Sym("squaref", &squaref); err != nil {
		t.Fatal(err)
	}

	if r := square(4); r != 16 {
		t.Errorf("expecting square(4) = 16, got %v instead", r)
	}
	if r := squaref(100); r != 10000 {
		t.Errorf("expecting squaref(100) = 10000, got %v instead", r)
	}
	var strlength func(string, string, string) int32
	if err := dl.Sym("strlength", &strlength); err != nil {
		t.Fatal(err)
	}
	const s = "this is not a long string"
	if r := int(strlength(s, s, s)); r != 3*len(s) {
		t.Errorf("expecting strlength(%q, %q, %q) = %v, got %v instead", s, s, s, 3*len(s), r)
	}

	var add func(int32, int32) int32
	if err := dl.Sym("add", &add); err != nil {
		t.Fatal(err)
	}
	if r := add(3, 2); r != 5 {
		t.Errorf("expecting add(3, 2) = 5, got %v instead", r)
	}

	var fill42 func([]byte, int32)
	if err := dl.Sym("fill42", &fill42); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 42)
	fill42(b, int32(len(b)))
	for ii, v := range b {
		if v != 42 {
			t.Errorf("b[%d] = %v != 42", ii, v)
		}
	}
}

func TestStackArguments(t *testing.T) {
	dl := openTestLib(t)
	defer dl.Close()
	var sum6 func(int32, int32, int32, int32, int32, int32) int32
	if err := dl.Sym("sum6", &sum6); err != nil {
		t.Fatal(err)
	}
	if r := sum6(1, 1, 1, 1, 1, 1); r != 6 {
		t.Errorf("expecting sum6(1...) = 6, got %v instead", r)
	}

	var sum8 func(int32, int32, int32, int32, int32, int32, int32, int32) int32
	if err := dl.Sym("sum8", &sum8); err != nil {
		t.Fatal(err)
	}
	if r := sum8(1, 2, 3, 4, 5, 6, 7, 8); r != 36 {
		t.Errorf("expecting sum8(1...8) = 36, got %v instead", r)
	}
}

func TestReturnString(t *testing.T) {
	dl := openTestLib(t)
	defer dl.Close()

	var returnString func(int32) string
	if err := dl.Sym("return_string", &returnString); err != nil {
		t.Fatal(err)
	}
	if r := returnString(0); r != "" {
		t.Errorf("expecting returnString(0) = \"\", got %v instead", r)
	}
	if r := returnString(1); r != "" {
		t.Errorf("expecting returnString(1) = \"\", got %v instead", r)
	}
	if r := returnString(2); r != "non-empty" {
		t.Errorf("expecting returnString(2) = \"non-empty\", got %v instead", r)
	}

	var returnStringPtr func(int32) *string
	if err := dl.Sym("return_string", &returnStringPtr); err != nil {
		t.Fatal(err)
	}
	if r := returnStringPtr(0); r != nil {
		t.Errorf("expecting returnStringPtr(0) = nil, got %v instead", r)
	}
	if r := returnStringPtr(1); r == nil || *r != "" {
		t.Errorf("expecting returnStringPtr(1) = \"\", got %v instead", r)
	}
	if r := returnStringPtr(2); r == nil || *r != "non-empty" {
		t.Errorf("expecting returnStringPtr(2) = \"non-empty\", got %v instead", r)
	}
}

func init() {
	if err := exec.Command("make", "-C", "testdata").Run(); err != nil {
		panic(err)
	}
}
