#ifndef VM_H
#define VM_H

#include "chunk.h"
#include "value.h"
#include "object.h"
#include "set.h"

#define STACK_MAX 4096
#define FRAMES_MAX 128
#define LOCALS_MAX 4096

struct call_frame {
    struct obj_closure *function;
    size_t old_ip;
    struct value *locals; // heap-allocated array
};

struct vm {
    struct chunk top_level;
    struct chunk *current;

    struct value *stack; // same
    struct value *stack_top;

    struct call_frame *frames; // same
    size_t frames_len;
    
    size_t ip;

    struct obj_upvalue *open_upvalues; // linked list, not owned by this field. owned by obj_list.
    struct object *obj_list;
    struct string_set strings;
};

struct call_frame call_frame_new(struct obj_closure *function, size_t old_ip);
void call_frame_free(struct call_frame *frame);

// ---

struct vm vm_new(struct chunk chunk, struct object *obj_list, struct string_set strings);
void vm_free(struct vm *vm);

bool vm_run(struct vm *vm);

#endif
