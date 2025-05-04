#ifndef VALUE_H
#define VALUE_H

typedef double float64;

typedef enum {
	VALUE_NUMBER,
	VALUE_BOOL,
	VALUE_NIL,
	VALUE_VOID,
} ValueType;

typedef struct {
	ValueType type;
	union {
		float64 number;
		bool boolean;
		
	} as;
} Value;

#define NEW_NUMBER(value) ((Value){VALUE_NUMBER, {.number = value}})
#define NEW_BOOL(value) ((Value){VALUE_BOOL, {.boolean = value}})

#endif
