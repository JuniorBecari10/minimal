// Parâmetro:
// INTERNAL_TYPE - Tipo

#include "util.h"

#include <stdlib.h>
#include <string.h>

#ifndef INTERNAL_TYPE
    #define INTERNAL_TYPE int
#endif

#define LIST_TYPE MACRO_CONCAT(List_, TEMPLATE_NAME)

#define INITIAL_CAP 10
#define LIST_GROWTH_FACTOR 2

struct LIST_TYPE {
    INTERNAL_TYPE* data;
    size_t size;
    size_t capacity;
};

static struct LIST_TYPE MACRO_CONCAT(LIST_TYPE, _new)(size_t size) {
	return (LIST_TYPE) {
		.data = (INTERNAL_TYPE *) malloc(sizeof(INTERNAL_TYPE) * INITIAL_CAP),
		.size = size,
		.capacity = INITIAL_CAP,
	};
}

static void MACRO_CONCAT(LIST_TYPE, _free)(struct LIST_TYPE *list) {
	free(list->data); // no problem if 'list->data' is NULL
    memset(list, 0, sizeof(*list));
}

static void MACRO_CONCAT(LIST_TYPE, _push)(struct LIST_TYPE *list, INTERNAL_TYPE value) {
    if (list->size + 1 > list->capacity) {
        list->capacity = (list->size + 1) * LIST_GROWTH_FACTOR;
        list->data = (INTERNAL_TYPE *) realloc(list->data, sizeof(*list->data) * list->capacity);
    }

    list->data[list->size++] = value;
}

static INTERNAL_TYPE MACRO_CONCAT(LIST_TYPE, _pop)(struct LIST_TYPE *list) {
    if (list->size > 0)
        return list->data[list->size--];

	// UNSAFE - don't use the value if the control path reaches here; always check the length first
	return 0;
}

static void MACRO_CONCAT(LIST_TYPE, _remove)(struct LIST_TYPE *list, size_t index) {
    if (index >= list->size)
        return;

    size_t amount_to_copy = list->size - 1 - index;

    if (amount_to_copy > 0)
        memmove(list->data + index, list->data + index + 1, amount_to_copy * sizeof(*list->data));
    
    list->size--;
}

#undef LIST_TYPE

