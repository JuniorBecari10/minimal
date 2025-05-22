#ifndef STRING_H
#define STRING_H

#include <inttypes.h>
#include <stddef.h>

struct string {
    char *chars;
    size_t length;

    uint32_t hash;
};

struct string string_new(const char *chars, size_t length);
struct string string_new_no_alloc(char *chars, size_t length);

void string_free(struct string *str);

#endif
