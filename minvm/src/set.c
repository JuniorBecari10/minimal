#include "set.h"

#include <stdlib.h>
#include <string.h>

static size_t hash(struct string_set *set, size_t key);
static void resize(struct string_set *set, size_t size);

struct string_set string_set_new(void) {
    return (struct string_set) {
        .strings = calloc(sizeof(struct string), INITIAL_CAPACITY),
        .length = 0,
        .capacity = INITIAL_CAPACITY,
    };
}

void string_set_free(struct string_set *set) {
    for (struct string *s = set->strings; s < set->strings + set->length; s++)
        string_free(s);

    free(set->strings);
    memset(set, 0, sizeof(*set));
}

struct string *string_set_add(struct string_set *set, struct string str) {
    if (set->length + 1 > set->capacity * LOAD_FACTOR)
        resize(set, set->capacity * GROW_FACTOR);
    
    size_t index = hash(set, (size_t) str.hash);

    for (struct string *s = set->strings + index;; index = (index + 1) % set->capacity) {
        if (s->chars == NULL) {
            // empty slot. insert here.
            *s = str;
            set->length++;

            return s;
        }

        else if (s->length == str.length && memcmp(s->chars, str.chars, s->length) == 0) {
            // value already exists. replace it.
            string_free(s);
            *s = str;

            return s;
        }
    }
}

struct string *string_set_get(struct string_set *set, struct string *str) {
    size_t index = hash(set, (size_t) str->hash);

    for (struct string *s = set->strings + index;; index = (index + 1) % set->capacity) {
        if (s->chars == NULL) {
            // empty slot. key not found.
            return NULL;
        }

        else if (s->length == str->length && memcmp(s->chars, str->chars, s->length) == 0) {
            // found the key. return it.
            return s;
        }
    }
}

// ---

static size_t hash(struct string_set *set, size_t key) {
    return key % set->capacity;
}

static void resize(struct string_set *set, size_t size) {
    set->strings = realloc(set->strings, size);
}
