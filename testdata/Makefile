all: lib.so

lib.so: lib.c
	gcc -shared -fPIC -o $@ $^
clean:
	$(RM) lib.so
