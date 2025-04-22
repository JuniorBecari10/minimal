#include <stdio.h>
#include <stdlib.h>

char *read_file(const char *path);

int main(int argc, char **argv) {
    if (argc != 2) {
        fprintf(stderr, "Usage: minvm <bytecode>");
    }

    char *file = read_file(argv[1]);

    free(file);
    return 0;
}

char *read_file(const char *path) {
    // TODO: read file
    return NULL;
}

