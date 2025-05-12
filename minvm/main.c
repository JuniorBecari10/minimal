#include "io.h"
#include <stdio.h>

#define TRY(e) if (!e) return 1

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: minvm <bytecode>\n");
        return 1;
    }
 
    const char *filename = argv[1];

    struct chunk chunk;
    struct object *obj_list;
    TRY(read_bytecode(filename, &chunk, &obj_list));

    return 0;
}
