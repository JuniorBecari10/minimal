#ifndef SET_H
#define SET_H

#include "string.h"

#define LOAD_FACTOR 0.75
#define INITIAL_CAPACITY 10
#define GROW_FACTOR 2

// for now, this will have a standalone implementation. if later we need maps, we include them here.
struct string_set {
    struct string *strings;
    size_t length;
    size_t capacity;
};

struct string_set string_set_new(void);
void string_set_free(struct string_set *set);

struct string *string_set_add(struct string_set *set, struct string str);
struct string *string_set_get(struct string_set *set, struct string *str);

#endif
