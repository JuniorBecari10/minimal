#include "io.h"
#include "object.h"
#include "set.h"
#include "vm.h"

#include <stdio.h>
#include <stdlib.h>

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: minvm <bytecode>\n");
        return 1;
    }
 
    char *filename = argv[1];

    struct chunk chunk = {0};
    struct object *obj_list = NULL;
    struct string_set strings = string_set_new();

    bool res = read_bytecode(filename, &chunk, &obj_list, &strings);

    // fill the open upvalues list when creating the VM.
    // the VM will take ownership of every argument passed to it.
    struct vm vm = vm_new(chunk, obj_list, strings);

    // reuse the free code in the vm.
    if (!res) {
        vm_free(&vm);
        return 1;
    }

    vm_run(&vm);
    vm_free(&vm);

    return 0;
}