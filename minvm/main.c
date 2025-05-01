#include <stdio.h>
#include <stdlib.h>

#include "include/util.h"
#include "include/io.h"

int main(int argc, char **argv) {
    if (argc != 2)
        ERROR_RET_1("Usage: minvm <bytecode>");

    char *file = read_file(argv[1]);
    printf("%s\n", file);
    
    free(file);
    return 0;
}

