#include "vm.h"

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

        .obj_list = obj_list,
        .strings = strings,  
    };
}

void vm_free(struct vm *vm) {
    chunk_free(&vm->top_level);
    string_set_free(&vm->strings);

    struct object *obj = vm->obj_list;
    while (obj != NULL) {
        object_free(obj);
        obj = obj->next;
    }
}

bool vm_run(struct vm *vm) {
    vm->current = &vm->top_level;
    vm->stack_top = vm->stack;

    return true;
}
