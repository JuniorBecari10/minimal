#ifndef UPVALUE_H
#define UPVALUE_H

#include "object.h"

typedef struct {
    Object obj;

    bool is_closed;
    union {
        Value closed;
        size_t index;
    } value;
} ObjUpvalue;

#define IS_UPVALUE(value) is_obj_type(value, OBJ_UPVALUE)
#define AS_UPVALUE(value) ((ObjUpvalue *) AS_OBJECT(value))

#endif
