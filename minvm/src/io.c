#include "../include/io.h"
#include "../include/util.h"

#include <stdio.h>
#include <stdlib.h>
#include <strings.h>
#include <stdbool.h>

// returns NULL if error.
char *read_file(const char *path, size_t *size) {
    FILE *file = NULL;

    if (strcasecmp(path, "*stdin") == 0)
        file = stdin;

    else {
        file = fopen(path, "rb");
        
        if (file == NULL)
            ERROR_RET_X("Cannot read file.", NULL);
    }

    size_t file_size = 0;
    char *buffer = NULL;

	if (file == stdin) {
        size_t capacity = 1024;
        size_t length = 0;
        
        buffer = (char *) malloc(capacity);
        
        if (!buffer)
            ERROR_RET_X("Memory allocation failed.", NULL);

        char c;
        while ((c = fgetc(file)) != EOF) {
            if (length + 1 >= capacity) {
                capacity *= 2;
                char *newBuffer = realloc(buffer, capacity);
                
                if (!newBuffer) {
                    free(buffer);
                    ERROR_RET_X("Memory allocation failed during read.", NULL);
                }

                buffer = newBuffer;
            }

            buffer[length++] = (char) c;
        }

        // Don't close stdin.
        
        buffer[length] = '\0';
        *size = length;
	} 
	
	else {
        fseek(file, 0L, SEEK_END);
        file_size = ftell(file);
        rewind(file);

        buffer = (char *) malloc(file_size + 1);
        
        if (!buffer)
            ERROR_RET_X("Memory allocation failed.", NULL);

        size_t bytes_read = fread(buffer, 1, file_size, file);
        buffer[bytes_read] = '\0';

        *size = bytes_read;
        fclose(file);
    }

    return buffer;
}

