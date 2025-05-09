#include "vm.h"
#include "object.h"
#include "set.h"
#include "instructions.h"

#include <stdlib.h>

VM init_vm(Chunk *chunk) {
	return (VM) {
		.chunk = chunk,
		.objects = NULL,
        .strings = string_set_new(),
	};
}

void free_vm(VM *vm) {
	// don't free chunk since it's not owned by the VM
    string_set_free(&vm->strings);

	Object *obj = vm->objects;
	while (obj != NULL) {
		free_object(obj);
		obj = obj->next;
	}
}

// this will only be freed through garbage collection or VM shutdown
void* allocate_object(VM* vm, size_t size, ObjType type) {
	// TODO: check and trigger GC if needed
	Object* object = (Object *) malloc(size);

	if (!object)
		return NULL;

	object->type = type;
	object->next = vm->objects;
	vm->objects = object;

	// TODO: increment bytes_allocated in VM (for garbage collection)
	return object;
}

bool interpret(VM *vm) {
	register uint8_t* ip = vm->chunk->code;

    #define NEXT_BYTE *(ip++)

    for (;;) {
        Instruction ins = (Instruction) NEXT_BYTE;

        switch (ins) {

        }
    }

    #undef NEXT_BYTE
}
