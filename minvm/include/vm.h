#ifndef VM_H
#define VM_H

#include "chunk.h"
#include "value.h"
#include "object.h"
#include "set.h"

#define STACK_MAX 4096
#define FRAMES_MAX 128

struct call_frame {
    struct obj_closure *function;

};

struct vm {
    struct chunk top_level;
    struct chunk *current;

    struct value stack[STACK_MAX];
    struct value *stack_top;

    struct call_frame frames[FRAMES_MAX];
    size_t frames_len;
    
    size_t ip;

    struct obj_upvalue *open_upvalues; // linked list, not owned by this field. owned by obj_list.
    struct object *obj_list;
    struct string_set strings;
};

struct vm vm_new(struct chunk chunk, struct object *obj_list, struct string_set strings);
void vm_free(struct vm *vm);

bool vm_run(struct vm *vm);

#endif
