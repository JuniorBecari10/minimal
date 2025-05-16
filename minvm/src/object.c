#include "object.h"

struct obj_string *new_string(const char *chars, size_t length) {

}

struct obj_function *new_function(struct chunk *chunk, size_t arity, char *name) {

}

struct obj_function *new_closure(struct obj_function *fn, struct obj_upvalue **upvalues, size_t upvalue_len) {

}

struct obj_native_fn *new_native_fn(native_fn *fn, size_t arity) {

}
