#ifndef OBJECT_H
#define OBJECT_H

#include "chunk.h"
#include "value.h"
#include "string.h"

#include <stdbool.h>
#include <stddef.h>

typedef struct value native_fn(struct value *args, size_t len);

enum object_type {
    OBJ_STRING,
    OBJ_FUNCTION,
    OBJ_CLOSURE,
    OBJ_NATIVE_FN,
    OBJ_UPVALUE,
};

struct object {
    enum object_type type;
    struct object *next;
};

// ---

struct obj_string {
    struct object obj;
    struct string *str;
};

struct obj_function {
    struct object obj;

    struct chunk chunk;
    size_t arity;
    char *name;
};

struct obj_closure {
    struct object obj;
    struct obj_function *fn;

    struct obj_upvalue **upvalues; // list to pointers to upvalues
    size_t upvalue_len;
};

struct obj_native_fn {
    struct object obj;

    size_t arity;
    native_fn *fn;
};

struct obj_upvalue {
    struct object obj;
    struct obj_upvalue *next; // intrusive list (part of vm's list)

    bool is_closed;

    union {
        struct value *location;
        struct value closed;
    } data;
};

struct object *object_new(size_t size, enum object_type type);
void object_free(struct object *obj);

static inline bool is_object_type(struct value value, enum object_type type) {
    return IS_OBJECT(value) && AS_OBJECT(value)->type == type;
}

struct obj_string *new_string(const char *chars, size_t length);
struct obj_function *new_function(struct chunk *chunk, size_t arity, char *name);
struct obj_function *new_closure(struct obj_function *fn, struct obj_upvalue **upvalues, size_t upvalue_len);
struct obj_native_fn *new_native_fn(native_fn *fn, size_t arity);
// TODO: new upvalue

#define IS_STRING(value)        is_object_type(value, OBJ_STRING)
#define AS_STRING(value)        ((struct obj_string *) AS_OBJECT(value))

#define IS_FUNCTION(value)      is_object_type(value, OBJ_FUNCTION)
#define AS_FUNCTION(value)      ((struct obj_function *) AS_OBJECT(value))

#define IS_CLOSURE(value)       is_object_type(value, OBJ_CLOSURE)
#define AS_CLOSURE(value)       ((struct obj_closure *) AS_OBJECT(value))

#define IS_NATIVE_FN(value)     is_object_type(value, OBJ_NATIVE_FN)
#define AS_NATIVE_FN(value)     ((struct obj_native_fn *) AS_OBJECT(value))

#define IS_UPVALUE(value)       is_object_type(value, OBJ_UPVALUE)
#define AS_UPVALUE(value)       ((struct obj_upvalue *) AS_OBJECT(value))

#endif
