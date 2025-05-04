#include "util.h"

#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#define DEFINE_LIST(INTERNAL_TYPE)                                                                          \
typedef struct List_##INTERNAL_TYPE {                                                                       \
    INTERNAL_TYPE* data;                                                                                    \
    size_t length;                                                                                          \
    size_t capacity;                                                                                        \
} List_##INTERNAL_TYPE;                                                                                     \
                                                                                                            \
static List_##INTERNAL_TYPE List_##INTERNAL_TYPE##_new() {                                                  \
    return (List_##INTERNAL_TYPE) {                                                                         \
        .data = (INTERNAL_TYPE *) malloc(sizeof(INTERNAL_TYPE) * 10),                                       \
        .length = 0,                                                                                        \
        .capacity = 10,                                                                                     \
    };                                                                                                      \
}                                                                                                           \
                                                                                                            \
static List_##INTERNAL_TYPE List_##INTERNAL_TYPE##_new_with_capacity(size_t capacity) {                     \
    return (List_##INTERNAL_TYPE) {                                                                         \
        .data = (INTERNAL_TYPE *) malloc(sizeof(INTERNAL_TYPE) * capacity),                                 \
        .length = 0,                                                                                        \
        .capacity = capacity,                                                                               \
    };                                                                                                      \
}                                                                                                           \
                                                                                                            \
static void List_##INTERNAL_TYPE##_free(List_##INTERNAL_TYPE *list) {                                       \
    free(list->data);                                                                                       \
    memset(list, 0, sizeof(*list));                                                                         \
}                                                                                                           \
                                                                                                            \
static bool List_##INTERNAL_TYPE##_push(List_##INTERNAL_TYPE *list, INTERNAL_TYPE value) {                  \
    if (list->length + 1 > list->capacity) {                                                                \
        size_t new_capacity = (list->length + 1) * 2;                                                       \
                                                                                                            \
        INTERNAL_TYPE* new_data = (INTERNAL_TYPE*) realloc(list->data, sizeof(*list->data) * new_capacity); \
        if (!new_data) return false;                                                                        \
                                                                                                            \
        list->data = new_data;                                                                              \
        list->capacity = new_capacity;                                                                      \
    }                                                                                                       \
                                                                                                            \
    list->data[list->length++] = value;                                                                     \
    return true;                                                                                            \
}                                                                                                           \
                                                                                                            \
static INTERNAL_TYPE List_##INTERNAL_TYPE##_pop(List_##INTERNAL_TYPE *list) {                               \
    if (list->length > 0)                                                                                   \
        return list->data[--list->length];                                                                  \
                                                                                                            \
    return (INTERNAL_TYPE){0};                                                                              \
}                                                                                                           \
                                                                                                            \
static void List_##INTERNAL_TYPE##_remove(List_##INTERNAL_TYPE *list, size_t index) {                       \
    if (index >= list->length) return;                                                                      \
                                                                                                            \
    size_t amount_to_copy = list->length - 1 - index;                                                       \
    if (amount_to_copy > 0)                                                                                 \
        memmove(&list->data[index], &list->data[index + 1], amount_to_copy * sizeof(*list->data));          \
                                                                                                            \
    list->length--;                                                                                         \
}                                                                                                           \
                                                                                                            \
static INTERNAL_TYPE List_##INTERNAL_TYPE##_get(const List_##INTERNAL_TYPE *list, size_t index) {           \
    if (index >= list->length)                                                                              \
        return (INTERNAL_TYPE){0};                                                                          \
                                                                                                            \
    return list->data[index];                                                                               \
}                                                                                                           \
                                                                                                            \
static void List_##INTERNAL_TYPE##_set(List_##INTERNAL_TYPE *list, size_t index, INTERNAL_TYPE value) {     \
    if (index < list->length)                                                                               \
        list->data[index] = value;                                                                          \
}
