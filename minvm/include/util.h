#ifndef UTIL_H
#define UTIL_H

#include <stdint.h>
#include <stddef.h>

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

#define MACRO_CONCAT(X, Y) X ## Y

uint32_t compute_checksum(uint8_t *data, size_t length);

#endif

