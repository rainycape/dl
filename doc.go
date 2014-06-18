// Package dl implements dynamic of shared libraries, like
// dlopen() and dlsym() in C.
//
//
// This package supports the following type mapping between Go
// and C.
//
//  Go		    -   C
//  (u)int	    -   (unsigned) long
//  (u)int8	    -   uint8_t / int8_t
//  (u)int16	    -   uint16_t / int16_t
//  (u)int32	    -   (unsigned) int
//  (u)int64	    -   uint64_t / int64_t
//  float32	    -	float
//  float64	    -	double
//  string	    -   char * (read only)
//  []byte	    -   char * (readwrite)
//  slices	    -   pointer to first argument
//  uintptr	    -   void *
//  unsafe.Pointer  -	void *
//
// No struct types are supported at this time.
//
// Retrieving variable symbols
//
// Symbols pointing to variables might be retrieved
// either as values or pointers. Given a C library
// which declares a symbol as:
//
//  int my_int = 8;
//
// It might be retrieved as a value, returning a copy with:
//
//  var myInt int32
//  if err := lib.Sym("my_int", &myInt); err != nil {
//	handle_error...
//  }
//
// Alternatively, a pointer to the variable might be obtained as:
//
//  var myInt *int32
//  if err := lib.Sym("my_int", &myInt); err != nil {
//	handle_error...
//  }
//
// Note that changing the value via the pointer will change the symbol
// in the loaded library, while changing the value obtained without the
// pointer will not, since a copy is made at lookup time.
//
// Retrieving function symbols
//
// This package also supports dynamically loading functions from libraries. To
// do so you must declare a function variable which matches the signature of the
// C function. Note that type mismatches will likely result in crashes, so use this
// feature with extreme care. Argument and return types must be of one of the
// supported types. See the examples in this package for the complete code.
//
//  var printf func(string, ...interface{}) int32
//  if err := lib.Sym("printf", &printf); err != nil {
//	handle_error...
//  }
//  printf("this string uses C format: %d\n", 7)
//
// Functions retrieved from a symbol can be used as standard Go functions.
//
// Overhead
//
// Typically, calling functions via this package rather than using cgo directly
// takes around 500ns more per call, due to reflection overhead. Future versions
// might adopt a JIT strategy which should make it as fast as cgo.
package dl
