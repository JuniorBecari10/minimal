#include "set.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

static size_t hash(struct string_set *set, size_t key);
static bool resize(struct string_set *set, size_t size);

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
}

struct string *string_set_add(struct string_set *set, struct string str) {
    if (set->length + 1 > set->capacity * LOAD_FACTOR && !resize(set, set->capacity * GROW_FACTOR))
        return NULL;

    size_t index = hash(set, str.hash);

    for (;; index = (index + 1) % set->capacity) {
        struct string *s = &set->strings[index];
        if (s->chars == NULL) {
            // Empty slot. Insert here.
            *s = str;
            set->length++;
            return s;
        }

        if (s->length == str.length && memcmp(s->chars, str.chars, s->length) == 0) {
            // Value already exists. Replace it.
            string_free(s);
            *s = str;
            return s;
        }
    }
}


struct string *string_set_get(struct string_set *set, struct string *str) {
    size_t index = hash(set, (size_t) str->hash);

    for (;; index = (index + 1) % set->capacity) {
        struct string *s = &set->strings[index];

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

static bool resize(struct string_set *set, size_t new_capacity) {
    struct string *old_strings = set->strings;
    size_t old_capacity = set->capacity;

    struct string *new_strings = calloc(new_capacity, sizeof(struct string));
    if (!new_strings)
        return false;

    set->strings = new_strings;
    set->capacity = new_capacity;
    set->length = 0;

    for (size_t i = 0; i < old_capacity; i++) {
        struct string *s = &old_strings[i];
        if (s->chars != NULL) {
            // Re-add to new table (ownership remains the same)
            string_set_add(set, *s);
        }
    }

    free(old_strings);
    return true;
}

