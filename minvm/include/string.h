#ifndef STRING_H
#define STRING_H

#include <stddef.h>

// the object that will be interned
// chars is heap-allocated and owned by the struct
typedef struct {
    char *chars;
    size_t length;
} String;

void string_free(String *s);

#endif
