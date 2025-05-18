#ifndef EXTRA_H
#define EXTRA_H

#include <stdint.h>

struct position {
    uint32_t line;
    uint32_t col;
};

struct metadata {
    struct position position;
    uint32_t length;
};

#endif
