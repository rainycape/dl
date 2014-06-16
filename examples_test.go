package dl_test

import (
	"bytes"
	"fmt"

	"dl"
)

func ExampleOpen_snprintf() {
	dl, err := dl.Open("libc", 0)
	if err != nil {
		panic(err)
	}
	var snprintf func([]byte, uint, string, ...interface{}) int
	if err := dl.Sym("snprintf", &snprintf); err != nil {
		panic(err)
	}
	buf := make([]byte, 200)
	snprintf(buf, uint(len(buf)), "hello %s!\n", "world")
	s := string(buf[:bytes.IndexByte(buf, 0)])
	fmt.Println(s)
	// Output: hello world!
}
