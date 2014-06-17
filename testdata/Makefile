CFLAGS=-O3

LIBS=

UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
	LIBS += lib.dylib
else
	LIBS += lib.so
endif

all: $(LIBS)

lib.so: lib.c
	gcc -shared -fPIC $(CFLAGS) -o $@ $^
lib.dylib: lib.c
	gcc -dynamiclib -undefined suppress -flat_namespace -fPIC $(CFLAGS) -o $@ $^
clean:
	$(RM) lib.so lib.dylib
