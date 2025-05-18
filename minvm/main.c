#include "io.h"
#include "object.h"
#include "set.h"

#include <stdio.h>
#include <stdlib.h>

#define TRY(e) if (!(e)) return 1

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: minvm <bytecode>\n");
        return 1;
    }
 
    char *filename = argv[1];

    struct chunk chunk = {0};
    struct object *obj_list = NULL;
    struct string_set strings = string_set_new();

    TRY(read_bytecode(filename, &chunk, &obj_list, &strings));
    free(filename);

    // the VM will take ownership of every argument passed to it.
    // VM vm = vm_new(chunk, obj_list, strings);
    // vm_free(&vm);

    return 0;
}
