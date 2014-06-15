package dl

import (
	"path/filepath"
	"testing"
	"unsafe"
)

func TestOpen(t *testing.T) {
	dl, err := Open("libc.so.6", 0)
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
	p := filepath.Join("testdata", "lib.so")
	dl, err := Open(p, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer dl.Close()
	var (
		s    string
		c    byte
		u32  uint32
		i    int
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
	p := filepath.Join("testdata", "lib.so")
	dl, err := Open(p, 0)
	if err != nil {
		t.Fatal(err)
	}
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
	var counter int
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
	dl, err := Open("libc.so.6", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer dl.Close()
	const s = "golang"
	var strlen func(string) int
	if err := dl.Sym("strlen", &strlen); err != nil {
		t.Error(err)
	} else if strlen(s) != len(s) {
		t.Errorf("expecting strlen(%q) = %v, got %v instead", s, len(s), strlen(s))
	}
}

func TestFunctions(t *testing.T) {
	p := filepath.Join("testdata", "lib.so")
	dl, err := Open(p, 0)
	if err != nil {
		t.Fatal(err)
	}
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
	var strlength func(string, string, string) int
	if err := dl.Sym("strlength", &strlength); err != nil {
		t.Fatal(err)
	}
	const s = "this is not a long string"
	if r := strlength(s, s, s); r != 3*len(s) {
		t.Errorf("expecting strlength(%q, %q, %q) = %v, got %v instead", s, s, s, 3*len(s), r)
	}

	var add func(int, int) int
	if err := dl.Sym("add", &add); err != nil {
		t.Fatal(err)
	}
	if r := add(3, 2); r != 5 {
		t.Errorf("expecting add(3, 2) = 5, got %v instead", r)
	}

	var fill42 func([]byte, int)
	if err := dl.Sym("fill42", &fill42); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 42)
	fill42(b, len(b))
	for ii, v := range b {
		if v != 42 {
			t.Errorf("b[%d] = %v != 42", ii, v)
		}
	}
}
