#ifndef OBJECT_H
#define OBJECT_H

#include "chunk.h"
#include "value.h"
#include "string.h"
#include "set.h"

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

    char *name;
    struct chunk chunk;
    size_t arity;
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
    struct obj_upvalue *next; // intrusive list (part of vm's open upvalues list)

    size_t upvalue_index;

    bool is_closed;
    union {
        struct value closed;
        struct value *location;
    } data;
};

// TODO: add a enum range_type field to indicate the range's type (int, float, char)
struct obj_range {
    struct object obj;

    struct value start;
    struct value end;
    struct value step;

    bool inclusive;
};

struct obj_record {
    struct object obj;
    char *name;

    char **field_names; // array of strings;
    size_t field_names_len;

    struct obj_closure **methods; // array of heap-allocated closures
    size_t methods_len;
};

struct obj_instance {
    struct object obj;
    struct obj_record *record; // pointer (borrow) to the record who created this instance
    
    struct value *fields;
    size_t fields_len;
};

struct obj_bound_method {
    struct object obj;

    struct value receiver;
    struct obj_closure *method;
};

struct object *object_new(size_t size, enum object_type type);
void object_free(struct object *obj);

void add_object_to_list(struct object *obj, struct object **list);
struct string *intern_string(struct string_set *set, struct string str);

static inline bool is_object_type(struct value value, enum object_type type) {
    return IS_OBJECT(value) && AS_OBJECT(value)->type == type;
}

struct obj_string *obj_string_new(struct string *str);
struct obj_function *obj_function_new(struct chunk chunk, size_t arity, char *name);
struct obj_closure *obj_closure_new(struct obj_function *fn, struct obj_upvalue **upvalues, size_t upvalue_len);
struct obj_native_fn *obj_native_fn_new(native_fn *fn, size_t arity);
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
