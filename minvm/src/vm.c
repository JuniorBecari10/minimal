#include "vm.h"
#include "object.h"
#include "instructions.h"
#include "value.h"

#include <stdint.h>
#include <stdlib.h>

static inline uint8_t next_byte(VM *vm);
static inline uint32_t next_uint32(VM *vm);
static inline void push(VM *vm, Value e);

VM init_vm(Chunk *chunk) {
    VM vm = {
		.chunk = chunk,
		.objects = NULL,

        .strings = string_map_new(),
        .stack = { {0} },
	};

    vm.ip = vm.chunk->code;
    vm.sp = vm.stack;

    return vm;
}

void free_vm(VM *vm) {
	// don't free chunk since it's not owned by the VM
    string_map_free(&vm->strings);

	Object *obj = vm->objects;
	while (obj != NULL) {
		object_free(obj);
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

String *intern_string(VM *vm, String str) {
    Entry *e = string_map_add(&vm->strings, str, NEW_NIL);
    return &e->key;
}

bool interpret(VM *vm) {
    for (;;) {
        Instruction ins = (Instruction) next_byte(vm);

        switch (ins) {
            case INST_PUSH_CONST: {
                uint32_t index = next_uint32(vm);
                Value constant = vm->chunk->constants.data[index];

                push(vm, constant);
            }
        }
    }
}

// ---

static inline uint8_t next_byte(VM *vm) {
    return *(vm->ip)++;
}

static inline uint32_t next_uint32(VM *vm) {
    return ((uint32_t)next_byte(vm))        |
           ((uint32_t)next_byte(vm) << 8)   |
           ((uint32_t)next_byte(vm) << 16)  |
           ((uint32_t)next_byte(vm) << 24);
}

static inline void push(VM *vm, Value e) {
    *vm->sp++ = e;
}
