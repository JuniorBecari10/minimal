#include "vm.h"
#include "instructions.h"

// these functions won't check bounds, since the caller has to uphold this guarantee.
uint8_t read_byte(struct vm *vm);

// these will.
bool push(struct vm *vm, struct value value);

struct vm vm_new(struct chunk chunk, struct object *obj_list, struct string_set strings) {
    // the current chunk will be set to top_level when running.
    return (struct vm) {
        .top_level = chunk,
        // current set in run

        .stack = { {0} },
        // stack_top set in run

        .frames = { {0} },
        .frames_len = 0,

        .ip = 0,

        .open_upvalues = NULL,
        .obj_list = obj_list,
        .strings = strings,
    };
}

void vm_free(struct vm *vm) {
    chunk_free(&vm->top_level);
    string_set_free(&vm->strings);

    for (struct object *obj = vm->obj_list; obj != NULL; obj = obj->next)
        object_free(obj);
}

bool vm_run(struct vm *vm) {
    vm->current = &vm->top_level;
    vm->stack_top = vm->stack;

    for (;;) {
        enum instruction inst = (enum instruction) read_byte(vm);

        switch (inst) {

        }
    }

    return true;
}

uint8_t read_byte(struct vm *vm) {
    return vm->current->code[vm->ip++];
}

bool push(struct vm *vm, struct value value) {

}
