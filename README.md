# About dl [![GoDoc][1]][2]

Runtime dynamic library loader (`dlopen`) for Go.

## Installing

Install in the usual Go fashion:

```sh
$ go get -u github.com/rainycape/dl
```

## Running Tests

`cgocheck` needs to be disabled in order to run the tests:

```sh
$ GODEBUG=cgocheck=0 go test -v
```

[1]: https://godoc.org/github.com/rainycape/dl?status.svg
[2]: https://godoc.org/github.com/rainycape/dl
