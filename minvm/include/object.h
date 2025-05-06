#ifndef OBJECT_H
#define OBJECT_H

#include "value.h"
#include "string.h"

#include <stdbool.h>
#include <stddef.h>

typedef enum {
    OBJ_STRING,
} ObjType;

typedef struct Object {
    ObjType type;
	struct Object *next;
} Object;

// the object that will be used as Value
typedef struct {
    Object obj;
	String *str;
} ObjString;

inline bool is_obj_type(Value value, ObjType type);
void free_object(Object *obj);

#define IS_STRING(value) is_obj_type(value, OBJ_STRING)

#define AS_STRING(value) ((ObjString *) AS_OBJECT(value))
#define AS_CSTRING(value) (((ObjString *) AS_OBJECT(value))->chars)

#endif
