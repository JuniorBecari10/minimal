#include "io.h"
#include "deserialize.h"
#include "object.h"

#include <stdio.h>
#include <stdlib.h>
#include <strings.h>

#define TRY(e) if (!e) return false

#define ERROR_RET(message, x)        \
    do {                               \
        fprintf(stderr, message "\n"); \
        return x;                      \
    } while (0)

#define HEADER "MNML"
#define HEADER_LEN 4

static char *read_file(const char *path, size_t *output_len);

static bool check_validity(const char *buffer, size_t len);

static bool check_header(const char *buffer, size_t file_len);
static bool check_checksum(const char *buffer, size_t file_len);

// TODO: use goto?
bool read_bytecode(const char *file_path, struct chunk *out, struct object **obj_list) {
    size_t buffer_len;
    char *buffer = read_file(file_path, &buffer_len);

    if (buffer == NULL)
        return false;

    if (!check_validity(buffer, buffer_len)) {
        free(buffer);
        return false;
    }

    bool res = deserialize(buffer, buffer_len, out, obj_list);
    free(buffer); // free the read file content unconditionally of the result of deserialize

    return res;
}


// returns NULL if error.
static char *read_file(const char *path, size_t *output_len) {
    FILE *file = NULL;

    if (strcasecmp(path, "*stdin") == 0)
        file = stdin;

    else {
        file = fopen(path, "rb");
        
        if (file == NULL)
            ERROR_RET("Cannot read file; file was not found.", NULL);
    }

    size_t file_size = 0;
    char *buffer = NULL;

	if (file == stdin) {
        size_t capacity = 1024;
        size_t length = 0;
        
        buffer = malloc(capacity);
        
        if (!buffer)
            ERROR_RET("Memory allocation failed.", NULL);

        char c;
        while ((c = fgetc(file)) != EOF) {
            if (length + 1 >= capacity) {
                capacity *= 2;
                char *newBuffer = realloc(buffer, capacity);
                
                if (!newBuffer) {
                    free(buffer);
                    ERROR_RET("Memory allocation failed during read.", NULL);
                }

                buffer = newBuffer;
            }

            buffer[length++] = (uint8_t) c;
        }

        // Don't close stdin.
        
        buffer[length] = '\0';
        *output_len = length;
	} 
	
	else {
        fseek(file, 0L, SEEK_END);
        file_size = ftell(file);
        rewind(file);

        buffer = malloc(file_size + 1);
        
        if (!buffer)
            ERROR_RET("Memory allocation failed.", NULL);

        size_t bytes_read = fread(buffer, 1, file_size, file);
        buffer[bytes_read] = '\0';

        *output_len = bytes_read;
        fclose(file);
    }

    return buffer;
}

static bool check_validity(const char *buffer, size_t len) {
    return check_header(buffer, len) && check_checksum(buffer, len);
}

static bool check_header(const char *buffer, size_t file_len) {
    
}

static bool check_checksum(const char *buffer, size_t file_len) {

}
