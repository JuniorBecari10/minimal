#include "object.h"
#include "set.h"
#include "value.h"
#include "chunk.h"
#include "string.h"

#include <stdlib.h>
#include <string.h>

#define TRY(e) if (!(e)) return NULL;

// size must be at least equal to the size of the desired object type.
// otherwise, it is UB.
// because of that, it is recommended to use the 'sizeof' operator in the 'size' parameter.
struct object *object_new(size_t size, enum object_type type) {
    struct object *obj = malloc(size);
    
    if (obj == NULL)
        return NULL;

    *obj = (struct object) {
        .type = type,
        .next = NULL,
    };

    return obj;
}

void object_free(struct object *obj) {
    switch (obj->type) {
        case OBJ_STRING: {
            // this does not free the string itself because it doesn't own it. the VM set does.
            struct obj_string *str = (struct obj_string *) obj;
            memset(str, 0, sizeof(*str));
            break;
        }

        case OBJ_FUNCTION: {
            struct obj_function *fn = (struct obj_function *) obj;
            
            free(fn->name);
            chunk_free(&fn->chunk);

            memset(fn, 0, sizeof(*fn));
            break;
        }

        case OBJ_CLOSURE: {
            struct obj_closure *closure = (struct obj_closure *) obj;
            object_free((struct object *) closure->fn);
            // TODO: free upvalues

            memset(closure, 0, sizeof(*closure));
            break;
        }

        case OBJ_UPVALUE: {
            struct obj_upvalue *upvalue = (struct obj_upvalue *) obj;
            
            // TODO: free upvalue.

            memset(upvalue, 0, sizeof(*upvalue));
            break;
        }

        // Native functions don't need to be freed.
        case OBJ_NATIVE_FN: break;
    }
}

// ---

void add_object_to_list(struct object *obj, struct object **list) {
    struct object *head = *list;

    obj->next = head;
    *list = obj;
}

struct string *intern_string(struct string_set *set, struct string str) {
    return string_set_add(set, str);
}

// ---

struct obj_string *obj_string_new(struct string *str) {
    struct obj_string *obj = (struct obj_string *) object_new(sizeof(struct obj_string), OBJ_STRING);
    TRY(obj);

    obj->str = str;
    return obj;
}

struct obj_function *obj_function_new(struct chunk chunk, size_t arity, char *name) {
    struct obj_function *obj = (struct obj_function *) object_new(sizeof(struct obj_function), OBJ_FUNCTION);
    TRY(obj);

    obj->chunk = chunk;
    obj->arity = arity;
    obj->name = name;

    return obj;
}

struct obj_closure *obj_closure_new(struct obj_function *fn, struct obj_upvalue **upvalues, size_t upvalue_len) {
    struct obj_closure *obj = (struct obj_closure *) object_new(sizeof(struct obj_closure), OBJ_CLOSURE);
    TRY(obj);

    obj->fn = fn;
    obj->upvalues = upvalues;
    obj->upvalue_len = upvalue_len;

    return obj;
}

struct obj_native_fn *obj_native_fn_new(native_fn *fn, size_t arity) {
    struct obj_native_fn *obj = (struct obj_native_fn *) object_new(sizeof(struct obj_native_fn), OBJ_NATIVE_FN);
    TRY(obj);

    obj->fn = fn;
    obj->arity = arity;

    return obj;
}
