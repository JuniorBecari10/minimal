#ifndef VM_H
#define VM_H

#include "chunk.h"
#include "object.h"

typedef struct {
	Chunk *chunk; // doesn't own
	Object *objects; // own
} VM;

VM init_vm(Chunk *chunk);
void free_vm(VM *vm);

void* allocate_object(VM* vm, size_t size, ObjType type);

#endif
