package dl

import (
	"C"
	"path/filepath"
	"runtime"
	"sync"
	"unsafe"
)

// Lib represents an opened dynamic library.
//
// Use Open to initialize a DL and use Close when finished making calls to the
// opened library.
//
// Note that when the DL is closed all its loaded symbols become invalid.
type Lib struct {
	h unsafe.Pointer
	sync.Mutex
}

// Open opens the shared library name with the given RTLD_* flags.
//
// If the library name does not have an extension, the platform default
// extension is appended to it (ie, .so, .dylib).
//
// See man dlopen for the meanings of the available RTLD_* flags.
//
// Note: if neither RTLD_LAZY nor RTLD_NOW are specified, then flags will be
// RTLD_NOW.
func Open(name string, flags int) (*Lib, error) {
	if flags&RTLD_LAZY == 0 && flags&RTLD_NOW == 0 {
		flags |= RTLD_NOW
	}

	// add extension
	if name != "" && filepath.Ext(name) == "" {
		name = name + LibExt
	}

	// open
	h, err := dlopen(name, flags)
	switch {
	case err != nil && runtime.GOOS == "linux" && name == "libc.so":
		// modern distros with libc6 have libc.so as a text file -- reattempt
		// opening with .6 suffix
		return Open(name+".6", flags)
	case err != nil:
		return nil, err
	}

	return &Lib{h: h}, nil
}

// Close closes the open library handle.
//
// Note: all previously loaded symbols will be invalidated.
func (l *Lib) Close() error {
	if l.h != nil {
		l.Lock()
		defer l.Unlock()
		defer func() { l.h = nil }()
		return dlclose(l.h)
	}
	return nil
}

// Sym loads symbol name from the loaded library into v.
func (l *Lib) Sym(name string, v interface{}) error {
	return Sym(l.h, name, v)
}
