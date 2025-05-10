#ifndef VM_H
#define VM_H

#include "chunk.h"
#include "object.h"
#include "map.h"
#include "string.h"
#include "value.h"

#define STACK_MAX 4096

typedef struct {
	Chunk *chunk; // doesn't own (currently executing Chunk)
	Object *objects; // own (currently live objects; not necessarily from the executing Chunk - it's a linked list)
    StringMap strings;
    Value stack[STACK_MAX];

    uint8_t *ip; // instruction pointer
    Value *sp;   // stack pointer
} VM;

VM init_vm(Chunk *chunk);
void free_vm(VM *vm);

void* allocate_object(VM* vm, size_t size, ObjType type);
String *intern_string(VM *vm, String str);

bool interpret(VM *vm);

#endif
