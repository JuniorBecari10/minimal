#ifndef SET_H
#define SET_H

#include "string.h"
#include "value.h"

#include <stdbool.h>

#define INITIAL_CAP 10

typedef struct {
    String key;
    Value value;
} Entry;

typedef struct {
	Entry *entries;
	size_t length;
	size_t capacity;
} StringMap;

StringMap string_map_new();
void string_map_free(StringMap *map);
void entry_free(Entry *entry);

Entry *string_map_add(StringMap* map, String str, Value value);
Entry *string_map_get(StringMap *map, String *str);

#endif

