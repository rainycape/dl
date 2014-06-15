
enum {
    ARG_FLAG_FLOAT = 1 << 0,
};

extern void *call(void *f, void **args, int *flags, int count);
