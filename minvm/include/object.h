#ifndef OBJECT_H
#define OBJECT_H

#include "value.h"

#include <stdbool.h>
#include <stddef.h>

typedef enum {
    OBJ_STRING,
} ObjType;

typedef struct Object {
    ObjType type;
    struct Object *next;
} Object;

inline bool is_obj_type(Value value, ObjType type);

typedef struct {
    Object obj;

    char *chars;
    size_t length;
} ObjString;

#define IS_STRING(value) is_obj_type(value, OBJ_STRING)

#define AS_STRING(value) ((ObjString *) AS_OBJECT(value))
#define AS_CSTRING(value) (((ObjString *) AS_OBJECT(value))->chars)

#endif
