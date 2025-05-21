#include "vm.h"
#include "instructions.h"
#include "object.h"
#include "error.h"

#include <stdlib.h>

#define TRY(e) if (!(e)) return false

// these functions won't check bounds, since the caller has to uphold this guarantee.
// and because these functions aren't supposed to fail.
static uint8_t read_byte(struct vm *vm);
static uint32_t read_uint32(struct vm *vm);

static struct obj_upvalue *capture_upvalue(struct vm *vm, struct value *location, size_t upvalue_index);

// these will.
static bool push(struct vm *vm, struct value value);

// operations
static bool push_const(struct vm *vm);
static bool push_closure(struct vm *vm);

// ---

struct call_frame call_frame_new(struct obj_closure *function, size_t old_ip) {
    return (struct call_frame) {
        .function = function,
        .old_ip = old_ip,
        .locals = malloc(LOCALS_MAX * sizeof(struct value)),
    };
}

void call_frame_free(struct call_frame *frame) {
    free(frame->locals);
}

// ---

struct vm vm_new(struct chunk chunk, struct object *obj_list, struct string_set strings) {
    // the current chunk will be set to top_level when running.
    return (struct vm) {
        .top_level = chunk,
        // current set in run

        .stack = malloc(STACK_MAX * sizeof(struct value)),
        // stack_top set in run

        .frames = malloc(STACK_MAX * sizeof(struct call_frame)),
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

    // expected to be empty when the vm is shut down, so no need to free the values here.
    free(vm->stack);

    for (struct call_frame *frame = vm->frames; frame < vm->frames + vm->frames_len; frame++)
        call_frame_free(frame);

    free(vm->frames);
}

bool vm_run(struct vm *vm) {
    vm->current = &vm->top_level;
    vm->stack_top = vm->stack;

    for (;;) {
        enum instruction inst = (enum instruction) read_byte(vm);

        switch (inst) {
            case INST_PUSH_CONST: TRY(push_const(vm)); break;
            case INST_PUSH_CLOSURE: TRY(push_closure(vm)); break;
        }
    }

    return true;
}

static uint8_t read_byte(struct vm *vm) {
    return vm->current->code[vm->ip++];
}

static uint32_t read_uint32(struct vm *vm) {
    uint32_t result =
        ((uint32_t) vm->current->code[vm->ip])           |
        ((uint32_t) vm->current->code[vm->ip + 1] << 8)  |
        ((uint32_t) vm->current->code[vm->ip + 2] << 16) |
        ((uint32_t) vm->current->code[vm->ip + 3] << 24);

    vm->ip += 4;
    return result;
}

// returns an upvalue to the specified variable.
// this function tries to find an existing one, but if it doesn't find one, it creates.
static struct obj_upvalue *capture_upvalue(struct vm *vm, struct value *location, size_t upvalue_index) {
    // try to find an existing upvalue to this variable.
    for (struct obj_upvalue *upvalue = vm->open_upvalues; upvalue != NULL; upvalue = upvalue->next) {
        if (upvalue->data.location == location && upvalue->upvalue_index == upvalue_index)
            return upvalue;
    }

    // not found. create a new one.
    struct obj_upvalue *upvalue = obj_upvalue_new_open(location, upvalue_index);
    upvalue->next = vm->open_upvalues;
    vm->open_upvalues = upvalue;

    return upvalue;
}

// ---

static bool push(struct vm *vm, struct value value) {
    if (vm->stack_top + 1 > vm->stack + STACK_MAX) {
        print_error(vm, "Operation stack exceeded.");
        return false;
    }

    *vm->stack_top++ = value;
    return true;
}

// ---

static bool push_const(struct vm *vm) {
    const struct value constant = vm->current->constants[read_byte(vm)];
    return push(vm, constant);
}

static bool push_closure(struct vm *vm) {
    struct obj_function *fn = AS_FUNCTION(vm->current->constants[read_byte(vm)]);
    const uint32_t upvalue_len = read_uint32(vm);

    // if there are no upvalues there's no need to allocate, and free can be unconditional since free(NULL) does nothing.
    struct obj_upvalue **upvalues = upvalue_len > 0
        ? malloc(upvalue_len * sizeof(struct obj_upvalue))
        : NULL;

    for (struct obj_upvalue **upvalue = upvalues; upvalue < upvalues + upvalue_len; upvalue++) {
        const bool is_local = read_byte(vm);
        const uint32_t index = read_uint32(vm);

        // if it's a local, create an upvalue and put it there.
        if (is_local)
            *upvalue = capture_upvalue(vm, vm->frames_len - 1, index);
        
        // if it's not, get it from the enclosing function's upvalue list.
        // we can safely get the last call frame since global functions do not capture upvalues.
        else
            *upvalue = vm->frames[vm->frames_len - 1].function->upvalues[index];
    }

    return push(vm, NEW_OBJECT(obj_closure_new(fn, upvalues, upvalue_len)));
}
