#include <stdio.h>
#include <stdlib.h>

#include "include/util.h"
#include "include/io.h"
#include "include/deserialize.h"

int main(int argc, char **argv) {
    if (argc != 2)
        ERROR_RET_1("Usage: minvm <bytecode>");

    size_t len;
    uint8_t *buffer = read_file(argv[1], &len);

	if (!check_validity(buffer, len)) {
		free(buffer);
		ERROR_RET_1("Invalid bytecode file.");
	}

	Chunk out;
	if (!deserialize(buffer, len, &out)) {
		free(buffer);
		ERROR_RET_1("Cannot read file.");
	}
	
    free(buffer);
    return 0;
}
