all: lib.so

lib.so: lib.c
	gcc -shared -fPIC -O3 -o $@ $^
clean:
	$(RM) lib.so
