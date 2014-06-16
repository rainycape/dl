#include <string.h>

#include "trampoline.h"

#define MAX_STACK_COUNT 100
#define MAX_INTEGER_COUNT (6)
#define MAX_FLOAT_COUNT (8)

#define _xstr(s) _str(s)
#define _str(s) #s

extern void * make_call(void *fn, void *regs, void *floats, int stack_count, void *stack, int is_float);

int
call(void *f, void **args, int *flags, int count, void **out)
{
    void *integers[MAX_INTEGER_COUNT];
    void *floats[MAX_FLOAT_COUNT];
    void *stack[MAX_STACK_COUNT];
    int integer_count = 0;
    int float_count = 0;
    int stack_count = 0;
    int ii;
    for (ii = 0; ii < count; ii++) {
        if (flags[ii] & ARG_FLAG_FLOAT) {
            if (float_count < MAX_FLOAT_COUNT) {
                floats[float_count++] = args[ii];
                continue;
            }
        } else {
            if (integer_count < MAX_INTEGER_COUNT) {
                integers[integer_count++] = args[ii];
                continue;
            }
        }
        if (stack_count > MAX_STACK_COUNT) {
            *out = strdup("maximum number of stack arguments reached (" _xstr(MAX_STACK_COUNT) ")");
            return 1;
        }
        // Argument on the stack
        stack[stack_count++] = args[ii];
    }
    void *floats_ptr = NULL;
    if (float_count > 0) {
        floats_ptr = floats;
    }
    if (stack_count & 1) {
        stack_count++;
    }
    for (ii = 0; ii < stack_count / 2; ii++) {
        int idx = stack_count-1-ii;
        void *tmp = stack[idx];
        stack[idx] = stack[ii];
        stack[ii] = tmp;
    }
    *out = make_call(f, integers, floats_ptr, stack_count, stack, flags[count] & ARG_FLAG_FLOAT);
    return 0;
}
