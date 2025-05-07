#ifndef VM_H
#define VM_H

#include "chunk.h"
#include "object.h"

typedef struct {
	Chunk *chunk; // doesn't own (currently executing Chunk)
	Object *objects; // own (currently live objects; not necessarily from the executing Chunk)
} VM;

VM init_vm(Chunk *chunk);
void free_vm(VM *vm);

void* allocate_object(VM* vm, size_t size, ObjType type);
bool interpret(VM *vm);

#endif
