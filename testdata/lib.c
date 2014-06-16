#include <stdint.h>
#include <string.h>
#include <stdio.h>

const char *my_string = "mystring";

char my_char = 42;

uint32_t my_uint32 = 1337;

int my_int = -9000;

void *my_pointer = (void *)0xdeadbeef;

int counter = 0;

void
increase_counter(void)
{
    counter++;
}

double
square(double a)
{
    printf("square(%f)\n", a);
    return a * a;
}

float
squaref(float a)
{
    printf("squaref(%f)\n", a);
    return a * a;
}

int
strlength(const char *s1, const char *s2, const char *s3)
{
    printf("strlength(%s, %s, %s)\n", s1, s2, s3);
    return strlen(s1) + strlen(s2) + strlen(s3);
}

int
add(int a, int b)
{
    printf("add(%d, %d)\n", a, b);
    return a + b;
}

void
fill42(unsigned char *data, int count) {
    int ii;
    printf("fill42 %p %d\n", data, count);
    for (ii = 0; ii < count; ii++) {
        data[ii] = 42;
    }
}

int
sum6(int a1, int a2, int a3, int a4, int a5, int a6)
{
    printf("sum6(%d, %d, %d, %d, %d, %d)\n", a1, a2, a3, a4, a5, a6);
    return a1 + a2 + a3 + a4 + a5 + a6;
}

int
sum8(int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8)
{
    //printf("sum8(%d, %d, %d, %d, %d, %d, %d, %d)\n", a1, a2, a3, a4, a5, a6, a7, a8);
    return a1 + a2 + a3 + a4 + a5 + a6 + a7 + a8;
}

int
ret8(int a1, int a2, int a3, int a4, int a5, int a6, int a7, int a8)
{
    printf("ret8(%d, %d, %d, %d, %d, %d, %d, %d)\n", a1, a2, a3, a4, a5, a6, a7, a8);
    return a8;
}
