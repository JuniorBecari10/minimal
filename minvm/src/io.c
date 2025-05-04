#include "../include/io.h"
#include "../include/util.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <stdbool.h>

// returns NULL if error.
uint8_t *read_file(const char *path, size_t *output_len) {
    FILE *file = NULL;

    if (strcasecmp(path, "*stdin") == 0)
        file = stdin;

    else {
        file = fopen(path, "rb");
        
        if (file == NULL)
            ERROR_RET_X("Cannot read file.", NULL);
    }

    size_t file_size = 0;
    uint8_t *buffer = NULL;

	if (file == stdin) {
        size_t capacity = 1024;
        size_t length = 0;
        
        buffer = (uint8_t *) malloc(capacity);
        
        if (!buffer)
            ERROR_RET_X("Memory allocation failed.", NULL);

        char c;
        while ((c = fgetc(file)) != EOF) {
            if (length + 1 >= capacity) {
                capacity *= 2;
                uint8_t *newBuffer = realloc(buffer, capacity);
                
                if (!newBuffer) {
                    free(buffer);
                    ERROR_RET_X("Memory allocation failed during read.", NULL);
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

        buffer = (uint8_t *) malloc(file_size + 1);
        
        if (!buffer)
            ERROR_RET_X("Memory allocation failed.", NULL);

        size_t bytes_read = fread(buffer, 1, file_size, file);
        buffer[bytes_read] = '\0';

        *output_len = bytes_read;
        fclose(file);
    }

    return buffer;
}

bool check_validity(const uint8_t *buffer, size_t len) {
	uint32_t checksum = compute_checksum(buffer, len - HEADER_LEN);
	uint8_t checksum_bytes[4] = {
		checksum       & 0xFF,
		checksum >> 8  & 0xFF,
		checksum >> 16 & 0xFF,
		checksum >> 24 & 0xFF,
	};

	return len > HEADER_LEN + CHECKSUM_LEN
		&& strncmp((char *) buffer, HEADER, HEADER_LEN) == 0
		&& strncmp(((char *) buffer) + len - CHECKSUM_LEN, (char *) checksum_bytes, CHECKSUM_LEN) == 0;
}
