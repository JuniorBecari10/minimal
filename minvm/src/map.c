#include "map.h"

#include <stdlib.h>
#include <stdbool.h>
#include <string.h>
#include <stddef.h>
#include <inttypes.h>

#define INITIAL_CAP 10
#define LOAD_FACTOR 0.75
#define GROW_FACTOR 2

static void string_map_resize(StringMap *map, size_t size);

static size_t hash(size_t key, size_t capacity);
static uint32_t hash_string(String *s);

StringMap string_map_new() {
	return (StringMap) {
		.entries = calloc(sizeof(Entry), INITIAL_CAP),
        .length = 0,
        .capacity = INITIAL_CAP,
	};
}

void string_map_free(StringMap *map) {
    size_t final_len = (size_t) map->entries + map->length;
    for (Entry *e = map->entries; (size_t) e < final_len; e++)
        entry_free(e);

    free(map->entries);
    memset(map, 0, sizeof(*map));
}

void entry_free(Entry *entry) {
    string_free(&entry->key);
    value_free(&entry->value);

    memset(entry, 0, sizeof(*entry));
}

// this takes ownership of str and value
Entry *string_map_add(StringMap* map, String str, Value value) {
    if (map->length + 1 > map->capacity * LOAD_FACTOR)
        string_map_resize(map, map->capacity * GROW_FACTOR);
    
    size_t hash_str = hash_string(&str);
    size_t index = hash(hash_str, map->capacity);

    for (
        Entry *entry = &map->entries[index];;
        index = (index + 1) % map->capacity
    ) {
        if (entry->key.chars == NULL) {
            // Empty slot; insert here
            entry->key = str;
            entry->value = value;

            map->length++;
            return entry;
        }

        else if (entry->key.length == str.length &&
                    memcmp(entry->key.chars, str.chars, str.length) == 0) {
            // Key already exists; update value and discard str
            value_free(&entry->value);
            string_free(&str);

            entry->value = value;
            return entry;
        }
    }
}

// borrows str
Entry *string_map_get(StringMap *map, String *str) {
    size_t hash_str = hash_string(str);
    size_t index = hash(hash_str, map->capacity);

    for (
        Entry *entry = &map->entries[index];;
        index = (index + 1) % map->capacity
    ) {
        if (entry->key.chars == NULL) {
            // Empty slot; key not found
            return NULL;
        }

        else if (entry->key.length == str->length &&
                 memcmp(entry->key.chars, str->chars, str->length) == 0) {
            // Found the entry; return it.
            return entry;
        }
    }
}

// 'size' must be greater than 'map->capacity'
static void string_map_resize(StringMap *map, size_t size) {
    map->entries = realloc(map->entries, size);
}

// ---

static size_t hash(size_t key, size_t capacity) {
    return key % capacity;
}

static uint32_t hash_string(String *s) {
    uint32_t hash = 2166136261u;
  
    for (size_t i = 0; i < s->length; i++) {
        hash ^= (uint8_t) s->chars[i];
        hash *= 16777619;
    }
    
    return hash;
}

