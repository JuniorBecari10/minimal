#ifndef VALUE_H
#define VALUE_H

#include <stdbool.h>

typedef double float64;

// defined in 'object.h'
struct object;

enum value_type {
    VALUE_NUMBER,
    VALUE_BOOL,
    VALUE_NIL,
    VALUE_VOID,
    VALUE_OBJ,
};

struct value {
    enum value_type type;

    union {
        float64 number;
        bool boolean;
        struct object *obj;
    } as;
};

#define IS_NUMBER(value)  ((value).type == VALUE_NUMBER)
#define IS_BOOL(value)    ((value).type == VALUE_BOOL)
#define IS_NIL(value)     ((value).type == VALUE_NIL)
#define IS_VOID(value)    ((value).type == VALUE_VOID)
#define IS_OBJECT(value)  ((value).type == VALUE_OBJ)

#define AS_NUMBER(value)  ((value).as.number)
#define AS_BOOL(value)    ((value).as.boolean)
#define AS_NIL NEW_NIL
#define AS_VOID NEW_VOID
#define AS_OBJECT(value)  ((value).as.obj)

#define NEW_NUMBER(value) ((struct value) { VALUE_NUMBER, { .number = value } })
#define NEW_BOOL(value)   ((struct value) { VALUE_BOOL, { .boolean = value } })
#define NEW_NIL(value)    ((struct value) { VALUE_NIL, { 0 } })
#define NEW_VOID(value)   ((struct value) { VALUE_VOID, { 0 } })
#define NEW_OBJECT(value) ((struct value) { VALUE_OBJ, { .obj = (struct object *) value } })

#endif
