#ifndef UTIL_H
#define UTIL_H

#define ERROR_RET_X(message, x)        \
    do {                               \
        fprintf(stderr, message "\n"); \
        return x;                      \
    } while (0)

#define ERROR_RET_1(message) ERROR_RET_X(message, 1)

#define ERRORF_RET_X(message, x, ...)               \
    do {                                            \
        fprintf(stderr, message "\n", __VA_ARGS__); \
        return x;                                   \
    } while (0)

#endif

