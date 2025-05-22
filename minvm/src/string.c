#include "string.h"

#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

static uint32_t hash_string(const char *chars, size_t length);

struct string string_new(const char *chars, size_t length) {
    // copy the string
    char *chars_heap = malloc(length + 1);

    strncpy(chars_heap, chars, length);
    chars_heap[length] = '\0';

    return string_new_no_alloc(chars_heap, length);
}

// constructs a string, but assumes that 'chars' is heap-allocated and valid.
struct string string_new_no_alloc(char *chars, size_t length) {
    return (struct string) {
        .chars = chars,
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
