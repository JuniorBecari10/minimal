#include "object.h"
#include "set.h"
#include "value.h"
#include "chunk.h"
#include "string.h"

#include <stdlib.h>

// size must be at least equal to the size of the desired object type.
// otherwise, it is UB.
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
            struct obj_string *str = (struct obj_string *) obj;
            string_free(str->str);

            break;
        }

        case OBJ_FUNCTION: {
            struct obj_function *fn = (struct obj_function *) obj;
            
            free(fn->name);
            chunk_free(&fn->chunk);

            break;
        }

        case OBJ_CLOSURE: {
            struct obj_closure *closure = (struct obj_closure *) obj;
            object_free((struct object *) closure->fn);
            // TODO: free upvalues

            break;
        }

        case OBJ_UPVALUE: {
            struct obj_upvalue *upvalue = (struct obj_upvalue *) obj;
            
            if (upvalue->is_closed && IS_OBJECT(upvalue->data.closed))
                object_free(AS_OBJECT(upvalue->data.closed));

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

struct string *intern_string(struct string str, struct string_set *set) {
    return string_set_add(set, str);
}

// ---

struct obj_string *obj_string_new(struct string *str) {

}

struct obj_function *obj_function_new(struct chunk chunk, size_t arity, char *name) {

}

struct obj_function *obj_closure_new(struct obj_function *fn, struct obj_upvalue **upvalues, size_t upvalue_len) {

}

struct obj_native_fn *obj_native_fn_new(native_fn *fn, size_t arity) {

}

