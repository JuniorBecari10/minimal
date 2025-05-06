#ifndef STRING_H
#define STRING_H

#include <stddef.h>

// the object that will be interned
typedef struct {
    char *chars;
    size_t length;
} String;

#endif
