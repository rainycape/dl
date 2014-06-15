package dlsym

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
