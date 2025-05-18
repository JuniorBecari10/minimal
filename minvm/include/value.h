#ifndef VALUE_H
#define VALUE_H

#include <stdbool.h>
#include <inttypes.h>

typedef double float64;

// defined in 'object.h'
struct object;

enum value_type {
    VALUE_INT,
    VALUE_FLOAT,
    VALUE_BOOL,
    VALUE_CHAR,
    VALUE_NIL,
    VALUE_VOID,
    VALUE_OBJ,
};

struct value {
    enum value_type type;

    union {
        int32_t integer;
        float64 floating;
        bool boolean;
        char character;
        struct object *obj;
    } as;
};

#define IS_INT(value)       ((value).type == VALUE_INT)
#define IS_FLOAT(value)     ((value).type == VALUE_FLOAT)
#define IS_BOOL(value)      ((value).type == VALUE_BOOL)
#define IS_CHAR(value)      ((value).type == VALUE_CHAR)
#define IS_NIL(value)       ((value).type == VALUE_NIL)
#define IS_VOID(value)      ((value).type == VALUE_VOID)
#define IS_OBJECT(value)    ((value).type == VALUE_OBJ)

#define AS_INT(value)       ((value).as.integer)
#define AS_FLOAT(value)     ((value).as.floating)
#define AS_BOOL(value)      ((value).as.boolean)
#define AS_CHAR(value)      ((value).as.character)
#define AS_NIL NEW_NIL
#define AS_VOID NEW_VOID
#define AS_OBJECT(value)    ((value).as.obj)

#define NEW_INT(_value)     ((struct value) { VALUE_INT, { .integer = _value } })
#define NEW_FLOAT(_value)   ((struct value) { VALUE_FLOAT, { .floating = _value } })
#define NEW_BOOL(_value)    ((struct value) { VALUE_BOOL, { .boolean = _value } })
#define NEW_CHAR(_value)    ((struct value) { VALUE_NUMBER, { .character = _value } })
#define NEW_NIL             ((struct value) { VALUE_NIL, { 0 } })
#define NEW_VOID            ((struct value) { VALUE_VOID, { 0 } })
#define NEW_OBJECT(_value)  ((struct value) { VALUE_OBJ, { .obj = (struct object *) _value } })

#endif
