#include "error.h"
#include "extra.h"

#include <stdio.h>

void print_error(struct vm *vm, const char *message) {
    struct metadata *meta = &vm->current->metadata[vm->ip];

    fprintf(stderr, "\n");
    fprintf(stderr, "[-] Runtime error: %s\n", message);
    fprintf(stderr, " |  [-] %s (%d, %d)\n", vm->current->name, meta->position.line + 1, meta->position.col + 1);
    fprintf(stderr, "[-]\n");
}
