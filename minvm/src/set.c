#include "set.h"

#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <stddef.h>
#include <inttypes.h>

#define INITIAL_CAP 10

static size_t hash(size_t key, size_t capacity);
static uint32_t hash_string(const char* key, size_t length);

StringSet string_set_new() {
	return (StringSet) {
		.entries = malloc(sizeof(String) * INITIAL_CAP),
        .length = 0,
        .capacity = INITIAL_CAP,
	};
}

void string_set_free(StringSet *set) {
    free(set->entries);
    memset(set, 0, sizeof(*set));
}

bool string_set_add(StringSet* set, String str) {

}

String string_set_get(StringSet *set, String str) {
    
}

// ---

static size_t hash(size_t key, size_t capacity) {
    return key % capacity;
}

static uint32_t hash_string(const char* key, size_t length) {
    uint32_t hash = 2166136261u;
  
    for (size_t i = 0; i < length; i++) {
        hash ^= (uint8_t) key[i];
        hash *= 16777619;
    }
    
    return hash;
}

