#ifndef SET_H
#define SET_H

#include "string.h"

#include <stdbool.h>

#define INITIAL_CAP 10

typedef struct {
	String *entries;
	size_t length;
	size_t capacity;
} StringSet;

StringSet string_set_new();
void string_set_free(StringSet *set);

bool string_set_add(StringSet* set, String str);
String string_set_get(StringSet *set, String str);

#endif

