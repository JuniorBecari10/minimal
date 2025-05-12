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
    OBJ_NATIVEFN,
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

    size_t arity;
    struct chunk chunk;
    char *name;
};

struct obj_closure {
    struct object obj;
    struct obj_function *fn;

    struct obj_upvalue **upvalues; // list to pointers to upvalues
    size_t upvalue_len;
};

struct obj_nativefn {
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

#endif
