#include "io.h"
#include "deserialize.h"
#include "object.h"
#include "checksum.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>

#define ERROR_RET(message, x)          \
    do {                               \
        fprintf(stderr, message "\n"); \
        return x;                      \
    } while (0)

static char *read_file(const char *path, size_t *output_len);

static bool check_validity(const char *buffer, size_t file_len);

static bool check_header(const char *buffer);
static bool check_checksum(const char *buffer, size_t file_len);

// TODO: use goto?
bool read_bytecode(const char *file_path,
                   struct chunk *out, struct object **obj_list, struct string_set *strings) {
    size_t buffer_len;
    char *buffer = read_file(file_path, &buffer_len);

    if (buffer == NULL)
        return false;

    if (!check_validity(buffer, buffer_len)) {
        fprintf(stderr, "Provided bytecode is not valid.\n");

        free(buffer);
        return false;
    }

    bool res = deserialize(buffer, buffer_len, out, obj_list, strings);
    free(buffer); // free the read file content unconditionally of the result of 'deserialize'.

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

static bool check_validity(const char *buffer, size_t file_len) {
    return file_len > HEADER_LEN + CHECKSUM_LEN &&
           check_header(buffer) &&
           check_checksum(buffer, file_len);
}

// 'len' is not needed, since we only need 4 bytes (defined by the macro)
// and 'check_validity' already checks it
static bool check_header(const char *buffer) {
    return strncmp(buffer, HEADER, HEADER_LEN) == 0;
}

static bool check_checksum(const char *buffer, size_t file_len) {
    uint32_t checksum = compute_checksum(buffer, file_len - HEADER_LEN);
	
    char checksum_bytes[4] = {
		checksum       & 0xFF,
		checksum >> 8  & 0xFF,
		checksum >> 16 & 0xFF,
		checksum >> 24 & 0xFF,
	};

    return strncmp(buffer + file_len - CHECKSUM_LEN, checksum_bytes, CHECKSUM_LEN) == 0;
}
