// Parameter:
// INTERNAL_TYPE - Type the list will hold

#include "util.h"

#include <stdlib.h>
#include <string.h>

#include <stdlib.h>
#include <string.h>

#define DEFINE_LIST(INTERNAL_TYPE)                                                                      \
typedef struct List_##INTERNAL_TYPE {                                                                   \
    INTERNAL_TYPE* data;                                                                                \
    size_t size;                                                                                        \
    size_t capacity;                                                                                    \
} List_##INTERNAL_TYPE;                                                                                 \
                                                                                                        \
static List_##INTERNAL_TYPE List_##INTERNAL_TYPE##_new(size_t size) {                                   \
    return (List_##INTERNAL_TYPE) {                                                                     \
        .data = (INTERNAL_TYPE *) malloc(sizeof(INTERNAL_TYPE) * 10),                                   \
        .size = size,                                                                                   \
        .capacity = 10,                                                                                 \
    };                                                                                                  \
}                                                                                                       \
                                                                                                        \
static void List_##INTERNAL_TYPE##_free(List_##INTERNAL_TYPE *list) {                                   \
    free(list->data);                                                                                   \
    memset(list, 0, sizeof(*list));                                                                     \
}                                                                                                       \
                                                                                                        \
static void List_##INTERNAL_TYPE##_push(List_##INTERNAL_TYPE *list, INTERNAL_TYPE value) {              \
    if (list->size + 1 > list->capacity) {                                                              \
        list->capacity = (list->size + 1) * 2;                                                          \
        list->data = (INTERNAL_TYPE *) realloc(list->data, sizeof(*list->data) * list->capacity);       \
    }                                                                                                   \
    list->data[list->size++] = value;                                                                   \
}                                                                                                       \
                                                                                                        \
static INTERNAL_TYPE List_##INTERNAL_TYPE##_pop(List_##INTERNAL_TYPE *list) {                           \
    if (list->size > 0)                                                                                 \
        return list->data[--list->size];                                                                \
    return (INTERNAL_TYPE){0};                                                                          \
}                                                                                                       \
                                                                                                        \
static void List_##INTERNAL_TYPE##_remove(List_##INTERNAL_TYPE *list, size_t index) {                   \
    if (index >= list->size)                                                                            \
        return;                                                                                         \
    size_t amount_to_copy = list->size - 1 - index;                                                     \
    if (amount_to_copy > 0)                                                                             \
        memmove(list->data + index, list->data + index + 1, amount_to_copy * sizeof(*list->data));      \
    list->size--;                                                                                       \
}                                                                                                       \
                                                                                                        \
static INTERNAL_TYPE List_##INTERNAL_TYPE##_get(const List_##INTERNAL_TYPE *list, size_t index) {       \
    if (index >= list->size)                                                                            \
        return (INTERNAL_TYPE){0};                                                                      \
    return list->data[index];                                                                           \
}                                                                                                       \
                                                                                                        \
static void List_##INTERNAL_TYPE##_set(List_##INTERNAL_TYPE *list, size_t index, INTERNAL_TYPE value) { \
    if (index < list->size)                                                                             \
        list->data[index] = value;                                                                      \
}

