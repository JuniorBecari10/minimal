#include <stdio.h>
#include <stdlib.h>

#include "include/util.h"

char *read_file(const char *path);

int main(int argc, char **argv) {
    if (argc != 2)
        ERROR_RET_1("Usage: minvm <bytecode>");

    char *file = read_file(argv[1]);

    printf("%s\n", file);
    free(file);
    return 0;
}

// returns NULL if error.
char *read_file(const char *path) {
    FILE *file = fopen(path, "rb");

    if (file == NULL)
        ERROR_RET_X("Cannot read file.", NULL);

    fseek(file, 0L, SEEK_END);
    size_t fileSize = ftell(file);
    rewind(file);
  
    char *buffer = (char *) malloc(fileSize + 1);
    size_t bytesRead = fread(buffer, sizeof(char), fileSize, file);
    buffer[bytesRead] = '\0';
  
    fclose(file);
    return buffer;
}

