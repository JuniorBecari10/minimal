#include "vm.h"
#include "instructions.h"
#include "object.h"
#include "error.h"
#include "options.h"
#include "value.h"

#include <math.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#define TRY(e) if (!(e)) return false

// these functions won't check bounds, since the caller has to uphold this guarantee.
// and because these functions aren't supposed to fail.
static uint8_t read_byte(struct vm *vm);
static uint32_t read_uint32(struct vm *vm);

static struct obj_upvalue *capture_upvalue(struct vm *vm, struct value *location, size_t upvalue_index);

// these will.
static bool push(struct vm *vm, struct value value);
static bool pop(struct vm *vm, struct value *out);

static bool assert_type(struct value value, enum value_type type);
static bool concatenate(struct vm *vm, const char *a, size_t a_len, const char *b, size_t b_len, char **out);

static struct object *vm_object_new(size_t size, enum object_type type);
static void vm_object_free(struct object *obj);

static void *vm_malloc(struct vm *vm, size_t size);

static bool binary(struct vm *vm, enum instruction operation);

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
            
            case INST_APPEND_METHODS: TRY((print_error(vm, "Uninmplemented"), false)); break;

            case INST_ADD:
            case INST_SUB:
            case INST_MUL:
            case INST_DIV:
            case INST_MOD: {
                TRY(binary(vm, inst));
                break;
            }
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
    struct obj_upvalue **upvalue_ptr = &vm->open_upvalues;

    // Find the correct place in the list (ordered by descending address).
    while (*upvalue_ptr != NULL && (*upvalue_ptr)->data.location > location)
        upvalue_ptr = &(*upvalue_ptr)->next;

    // Check if we already have an upvalue for this location and index.
    // It must be in this place since if it existed it would be here.
    if (*upvalue_ptr != NULL &&
        (*upvalue_ptr)->data.location == location &&
        (*upvalue_ptr)->upvalue_index == upvalue_index)
        return *upvalue_ptr;

    // Create a new upvalue and insert it into the list.
    // Inserting it here guarantees the descending order of the list.
    struct obj_upvalue *new_upvalue = obj_upvalue_new_open(location, upvalue_index);
    new_upvalue->next = *upvalue_ptr;
    *upvalue_ptr = new_upvalue;

    return new_upvalue;
}

// ---

static bool push(struct vm *vm, struct value value) {
    if (vm->stack_top + 1 >= vm->stack + STACK_MAX) {
        print_error(vm, "Operation stack exceeded.");
        return false;
    }

    *vm->stack_top++ = value;
    return true;
}

static bool pop(struct vm *vm, struct value *out) {
    if (vm->stack_top == vm->stack) {
        print_error(vm, "Attempt to pop an empty stack.");
        return false;
    }

    *out = *(--vm->stack_top);
    return true;
}

static bool assert_type(struct value value, enum value_type type) {
    return value.type == type;
}

// assumes length parameters do not contain '\0'.
static bool concatenate(struct vm *vm, const char *a, size_t a_len, const char *b, size_t b_len, char **out) {
    char *result = vm_malloc(vm, a_len + b_len + 1);
    TRY(result);

    memcpy(result, a, a_len);
    memcpy(result + a_len, b, b_len);

    result[a_len + b_len] = '\0';

    *out = result;
    return true;
}

// wrappers for allocate_object and free_object, but with bookmarking additions for the GC.
static struct object *vm_object_new(size_t size, enum object_type type) {
    // TODO: increment bytes_allocated
    return object_new(size, type);
}

static void vm_object_free(struct object *obj) {
    // TODO: decrement bytes_allocated
    object_free(obj);
}

static void *vm_malloc(struct vm *vm, size_t size) {
    // TODO: increment bytes_allocated
    return malloc(size);
}

static bool binary(struct vm *vm, enum instruction operation) {
    struct value left, right;

    TRY(pop(vm, &right));
    TRY(pop(vm, &left));

#ifdef ENABLE_TYPE_CHECKER
    if (left.type != right.type) {
        print_error(vm, "Types must be equal when performing arithmetic.");
        return false;
    }
#endif

    switch (operation) {
        case INST_ADD: {
            // int, float, str.

            switch (left.type) {
                case VALUE_INT: {
                    int32_t a = left.as.integer;
                    int32_t b = right.as.integer;

                    TRY(push(vm, NEW_INT(a + b)));
                    break;
                }
                
                case VALUE_FLOAT: {
                    float64 a = left.as.floating;
                    float64 b = right.as.floating;

                    TRY(push(vm, NEW_FLOAT(a + b)));
                    break;
                }
                
                case VALUE_OBJ: {
#ifdef ENABLE_TYPE_CHECKER
                    if (AS_OBJECT(left)->type != OBJ_STRING) {
                        print_error(vm, "Object type is not 'str'.");
                        return false;
                    }
#endif

                    struct obj_string *a = AS_STRING(left);
                    struct obj_string *b = AS_STRING(right);

                    char *concatenated;
                    TRY(concatenate(vm,
                                    AS_CSTRING(left), a->str->length,
                                    AS_CSTRING(right), b->str->length, &concatenated));

                    struct string str = string_new_no_alloc(concatenated, a->str->length + b->str->length);
                    struct string *interned = intern_string(&vm->strings, str);

                    struct obj_string *new_str =
                        (struct obj_string *) vm_object_new(sizeof(struct obj_string), OBJ_STRING);

                    TRY(new_str);
                    new_str->str = interned;

                    TRY(push(vm, NEW_OBJECT(new_str)));
                    break;
                }

                default: {
#ifdef ENABLE_TYPE_CHECKER
                    print_error(vm, "Invalid type for addition.");
                    return false;
#endif
                }
            }
        }
        
        case INST_SUB:
        case INST_MUL:
        case INST_DIV:
        case INST_MOD: {
            // int, float.
            
            switch (left.type) {
                case VALUE_INT: {
                    int32_t a = left.as.integer;
                    int32_t b = right.as.integer;

                    int32_t result = 0;
                    switch (operation) {
                        case INST_SUB: result = a - b; break;
                        case INST_MUL: result = a * b; break;
                        case INST_DIV: {
                            if (b == 0) {
                                print_error(vm, "Cannot divide by zero.");
                                return false;
                            }
                        
                            result = a / b;
                            break;
                        }
                        
                        case INST_MOD: {
                            if (b == 0) {
                                print_error(vm, "Cannot divide by zero.");
                                return false;
                            }
                        
                            result = a % b;
                            break;
                        }

                        default: { }
                    }

                    TRY(push(vm, NEW_INT(result)));
                    break;
                }
                
                case VALUE_FLOAT: {
                    float64 a = left.as.floating;
                    float64 b = right.as.floating;
                    
                    float64 result = 0;
                    switch (operation) {
                        case INST_SUB: result = a - b; break;
                        case INST_MUL: result = a * b; break;
                        case INST_DIV: {
                            if (b == 0.0) {
                                print_error(vm, "Cannot divide by zero.");
                                return false;
                            }

                            result = a / b;
                            break;
                        }
                        
                        case INST_MOD: {
                            if (b == 0.0) {
                                print_error(vm, "Cannot divide by zero.");
                                return false;
                            }
                        
                            result = fmod(a, b);
                            break;
                        }

                        default: { }
                    }

                    TRY(push(vm, NEW_FLOAT(result)));
                    break;
                }

                default: {
#ifdef ENABLE_TYPE_CHECKER
                    print_error(vm, "Invalid type for arithmetic.");
                    return false;
#endif
                }
            }
        }

        // cannot reach here because of the instruction switch.
        default: { }
    }

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
        // Get the address of the local in the current scope.
        // (globals are not captured, so this is safe)
        if (is_local)
            *upvalue = capture_upvalue(vm, &vm->frames[vm->frames_len - 1].locals[index], index);
        
        // if it's not, get it from the enclosing function's upvalue list.
        // we can safely get the last call frame since global functions do not capture upvalues.
        else
            *upvalue = vm->frames[vm->frames_len - 1].function->upvalues[index];
    }

    return push(vm, NEW_OBJECT(obj_closure_new(fn, upvalues, upvalue_len)));
}
