// Parameter:
// INTERNAL_TYPE - Type the list will hold

#include "util.h"

#include <stdlib.h>
#include <string.h>

#ifndef INTERNAL_TYPE
	#define INTERNAL_TYPE int
#endif

#define LIST_TYPE MACRO_CONCAT(List_, INTERNAL_TYPE)

#define INITIAL_CAP 10
#define LIST_GROWTH_FACTOR 2
#define DEFAULT_VALUE (INTERNAL_TYPE){0}

typedef struct LIST_TYPE {
    INTERNAL_TYPE* data;
    size_t size;
    size_t capacity;
} LIST_TYPE;

static LIST_TYPE MACRO_CONCAT(LIST_TYPE, _new)(size_t size) {
    return (LIST_TYPE) {
        .data = (INTERNAL_TYPE *) malloc(sizeof(INTERNAL_TYPE) * INITIAL_CAP),
        .size = size,
        .capacity = INITIAL_CAP,
    };
}

static void MACRO_CONCAT(LIST_TYPE, _free)(LIST_TYPE *list) {
    free(list->data);
    memset(list, 0, sizeof(*list));
}

static void MACRO_CONCAT(LIST_TYPE, _push)(LIST_TYPE *list, INTERNAL_TYPE value) {
    if (list->size + 1 > list->capacity) {
        list->capacity = (list->size + 1) * LIST_GROWTH_FACTOR;
        list->data = (INTERNAL_TYPE *) realloc(list->data, sizeof(*list->data) * list->capacity);
    }

    list->data[list->size++] = value;
}

static INTERNAL_TYPE MACRO_CONCAT(LIST_TYPE, _pop)(LIST_TYPE *list) {
    if (list->size > 0)
        return list->data[--list->size];

    return DEFAULT_VALUE;
}

static void MACRO_CONCAT(LIST_TYPE, _remove)(LIST_TYPE *list, size_t index) {
    if (index >= list->size)
        return;

    size_t amount_to_copy = list->size - 1 - index;

    if (amount_to_copy > 0)
        memmove(list->data + index, list->data + index + 1, amount_to_copy * sizeof(*list->data));

    list->size--;
}

static INTERNAL_TYPE MACRO_CONCAT(LIST_TYPE, _get)(const LIST_TYPE *list, size_t index) {
    if (index >= list->size)
        return DEFAULT_VALUE;

    return list->data[index];
}

static void MACRO_CONCAT(LIST_TYPE, _set)(LIST_TYPE *list, size_t index, INTERNAL_TYPE value) {
    if (index < list->size)
        list->data[index] = value;
}

#undef LIST_TYPE
#undef INTERNAL_TYPE

