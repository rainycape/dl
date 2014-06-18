#include <stdint.h>
#include <string.h>
#include <stdio.h>

const char *my_string = "mystring";

char my_char = 42;

uint32_t my_uint32 = 1337;

int my_int = -9000;

long my_long = -9000;

void *my_pointer = (void *)0xdeadbeef;

volatile int verbose = 0;

#define _printf(...) do { \
    if (verbose) { \
        printf(__VA_ARGS__); \
    } \
} while(0)

int counter = 0;

void
increase_counter(void)
{
    counter++;
}

double
square(double a)
{
    _printf("square(%f)\n", a);
    return a * a;
}

float
squaref(float a)
{
    _printf("squaref(%f)\n", a);
    return a * a;
}

int
strlength(const char *s1, const char *s2, const char *s3)
{
    _printf("strlength(%s, %s, %s)\n", s1, s2, s3);
    return strlen(s1) + strlen(s2) + strlen(s3);
}

int
add(int a, int b)
{
    _printf("add(%d, %d)\n", a, b);
    return a + b;
}

void
fill42(unsigned char *data, int count) {
    int ii;
    _printf("fill42 %p %d\n", data, count);
    for (ii = 0; ii < count; ii++) {
        data[ii] = 42;
    }
}

int
sum6(int a1, int a2, int a3, int a4, int a5, int a6)
{
    _printf("sum6(%d, %d, %d, %d, %d, %d)\n", a1, a2, a3, a4, a5, a6);
    return a1 + a2 + a3 + a4 + a5 + a6;
}

int
sum8(int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8)
{
    _printf("sum8(%d, %d, %d, %d, %d, %d, %d, %d)\n", a1, a2, a3, a4, a5, a6, a7, a8);
    return a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8;
}

int
ret8(int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8)
{
    _printf("ret8(%d, %d, %d, %d, %d, %d, %d, %d)\n", a1, a2, a3, a4, a5, a6, a7, a8);
    return a8;
}

char *
return_string(int a)
{
    if (a == 0) {
        return NULL;
    }
    if (a == 1) {
        return "";
    }
    return "non-empty";
}
