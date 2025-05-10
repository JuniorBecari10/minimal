#include "map.h"

#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <stddef.h>
#include <inttypes.h>

#define INITIAL_CAP 10

static size_t hash(size_t key, size_t capacity);
static uint32_t hash_string(const char* key, size_t length);

StringMap string_map_new() {
	return (StringMap) {
		.entries = malloc(sizeof(Entry) * INITIAL_CAP),
        .length = 0,
        .capacity = INITIAL_CAP,
	};
}

void string_map_free(StringMap *map) {
    free(map->entries);
    memset(map, 0, sizeof(*map));
}

Entry *string_map_add(StringMap* map, String str, Value value) {

}

Entry *string_map_get(StringMap *map, String str) {

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

