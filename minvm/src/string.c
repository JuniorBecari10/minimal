#include "string.h"

#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

static uint32_t hash_string(const char *chars, size_t length);

struct string string_new(const char *chars, size_t length) {
    char *chars_heap = malloc(length + 1);
    strncpy(chars_heap, chars, length);

    return (struct string) {
        .chars = chars_heap,
        .length = length,
        .hash = hash_string(chars, length),
    };
}

void string_free(struct string *str) {
    free(str->chars);
}

static uint32_t hash_string(const char *chars, size_t length) {
    uint32_t hash = 2166136261u;
  
    for (size_t i = 0; i < length; i++) {
        hash ^= (uint8_t) chars[i];
        hash *= 16777619;
    }
    
    return hash;
}
