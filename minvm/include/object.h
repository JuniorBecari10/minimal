#ifndef OBJECT_H
#define OBJECT_H

typedef enum {
    OBJ_STRING,
} ObjType;

typedef struct Object {
    ObjType type;
    struct Object *next;
} Object;

typedef struct {
    Object obj;
    
    char *chars;
    size_t length;
} ObjString;

#endif
