#include <stdio.h>
#include <stdlib.h>

#include "util.h"
#include "io.h"
#include "deserialize.h"

int main(int argc, char **argv) {
    if (argc != 2)
        ERROR_RET_1("Usage: minvm <bytecode>");

    size_t len;
    uint8_t *buffer = read_file(argv[1], &len);

	if (!buffer) return 1;

	if (!check_validity(buffer, len)) {
		free(buffer);
		ERROR_RET_1("Invalid bytecode file.");
	}

	Chunk out = {0}; // initialize everything to 0
	if (!deserialize(buffer, len, &out)) {
		free(buffer);
		free_chunk(&out);

		ERROR_RET_1("Cannot read file.");
	}
	
	free_chunk(&out);
    free(buffer);
    return 0;
}
