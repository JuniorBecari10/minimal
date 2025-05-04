#ifndef VALUE_H
#define VALUE_H

#include <stdbool.h>

typedef double float64;

// defined in 'object.h', but cannot import it because of cyclic references
struct Object;
typedef struct Object Object;

typedef enum {
	VALUE_NUMBER,
	VALUE_BOOL,
	VALUE_NIL,
	VALUE_VOID,
	VALUE_OBJ,
} ValueType;

typedef struct {
	ValueType type;
	
	union {
		float64 number;
		bool boolean;
		Object *obj;
	} as;
} Value;

#define IS_NUMBER(value) ((value).type == VALUE_NUMBER)
#define IS_BOOL(value) ((value).type == VALUE_BOOL)
#define IS_NIL(value) ((value).type == VALUE_NIL)
#define IS_VOID(value) ((value).type == VALUE_VOID)
#define IS_OBJ(value) ((value).type == VALUE_OBJ)

#define AS_NUMBER(value) ((value).as.number)
#define AS_BOOL(value) ((value).as.boolean)
#define AS_NIL NEW_NIL
#define AS_VOID NEW_VOID
#define AS_OBJECT(value) ((value).as.obj)

#define NEW_NUMBER(value) ((Value) {VALUE_NUMBER, {.number = value}})
#define NEW_BOOL(value) ((Value) {VALUE_BOOL, {.boolean = value}})
#define NEW_NIL ((Value) {VALUE_NIL, {0}})
#define NEW_VOID ((Value) {VALUE_VOID, {0}})
#define NEW_OBJECT(value) ((Value) {VALUE_OBJ, {.obj = (Object *) value}})

void free_value(Value *v);

#endif
