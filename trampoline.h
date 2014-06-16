
enum {
    ARG_FLAG_SIZE_8 = 1 << 0,
    ARG_FLAG_SIZE_16 = 1 << 1,
    ARG_FLAG_SIZE_32 = 1 << 2,
    ARG_FLAG_SIZE_64 = 1 << 3,
    ARG_FLAG_SIZE_PTR = 1 << 4,
    ARG_FLAG_FLOAT = 1 << 5,
};

extern int call(void *f, void **args, int *flags, int count, void **out);
