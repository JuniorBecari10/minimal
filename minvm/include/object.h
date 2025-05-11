#ifndef OBJECT_H
#define OBJECT_H

#include "chunk.h"
#include "value.h"
#include "upvalue.h"
#include "string.h"

#include <stdbool.h>
#include <stddef.h>

// use it like this: 'NativeFn *'.
typedef Value NativeFn(Value *args, size_t len);

typedef enum {
    OBJ_STRING,
    OBJ_FUNCTION,
} ObjType;

typedef struct Object {
    ObjType type;
	struct Object *next; // pointer to the next node of the VM's objects linked list.
} Object;

// the object that will be used as Value. it doesn't contain a string, but rather a pointer to the interned one.
typedef struct {
    Object obj;
	String *str;
} ObjString;

typedef struct {
    Object obj;

    size_t arity;
    Chunk chunk; // owns
    char *name; // may be NULL, since it's optional
} ObjFunction;

typedef struct {
    Object obj;

    ObjFunction *fn;
    ObjUpvalue **upvalues; // list of pointers to upvalues
} ObjClosure;

typedef struct {
    Object obj;

    size_t arity;
    NativeFn *fn;
} ObjNativeFn;

typedef struct {
    Object obj;

    float64 start;
    float64 end;
    float64 step;

    bool inclusive;
} ObjRange;

typedef struct {
    Object obj;

    ObjString name; // owned
    List_ObjString field_names;
    List_ObjClosure methods;
} ObjRecord;

// maybe use a map here?
typedef struct {
    Object obj;

    List_Value fields;
    ObjRecord *record;
} ObjInstance;

typedef struct {
    Object obj;

    Value receiver;
    ObjClosure method;
} ObjBoundMethod;

bool is_obj_type(Value value, ObjType type);
void object_free(Object *obj);

#define IS_STRING(value)        is_obj_type(value, OBJ_STRING)
#define AS_STRING(value)        ((ObjString *) AS_OBJECT(value))

#define AS_STRING_OBJ(value)    (AS_STRING(value)->str)
#define AS_CSTRING(value)       (AS_STRING(value)->str->chars)

#define IS_FUNCTION(value)      is_obj_type(value, OBJ_FUNCTION)
#define AS_FUNCTION(value)      ((ObjFunction *) AS_OBJECT(value))

#define IS_CLOSURE(value)       is_obj_type(value, OBJ_CLOSURE)
#define AS_CLOSURE(value)       ((ObjClosure *) AS_OBJECT(value))
#define CLOSURE_FN(value)       (AS_CLOSURE(value)->fn)

#define IS_NATIVE_FN(value)     is_obj_type(value, OBJ_NATIVE_FN)
#define AS_NATIVE_FN(value)     ((ObjNativeFn *) AS_OBJECT(value))

#define IS_RANGE(value)         is_obj_type(value, OBJ_RANGE)
#define AS_RANGE(value)         ((ObjRange *) AS_OBJECT(value))

#define IS_RECORD(value)        is_obj_type(value, OBJ_RECORD)
#define AS_RECORD(value)        ((ObjRecord *) AS_OBJECT(value))

#define IS_INSTANCE(value)      is_obj_type(value, OBJ_INSTANCE)
#define AS_INSTANCE(value)      ((ObjInstance *) AS_OBJECT(value))

#define IS_BOUND_METHOD(value)  is_obj_type(value, OBJ_BOUND_METHOD)
#define AS_BOUND_METHOD(value)  ((ObjBoundMethod *) AS_OBJECT(value))

#endif
