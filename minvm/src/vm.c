#include "../include/vm.h"
#include "../include/object.h"

#include <stdlib.h>

VM init_vm(Chunk *chunk) {
	return (VM) {
		.chunk = chunk,
		.objects = NULL,
	};
}

void free_vm(VM *vm) {
	// don't free chunk since it's not owned by the VM

	Object *obj = vm->objects;
	while (obj != NULL) {
		free_object(obj);
		obj = obj->next;
	}
}

void* allocate_object(VM* vm, size_t size, ObjType type) {
	// TODO: check and trigger GC if needed
	Object* object = (Object *) malloc(size);

	if (!object)
		return NULL;

	object->type = type;
	object->next = vm->objects;
	vm->objects = object;

	// TODO: increment bytes_allocated in VM
	return object;
}

bool interpret(VM *vm) {
	register uint8_t* ip = vm->chunk->code;
}
