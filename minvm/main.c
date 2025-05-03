#include <stdio.h>
#include <stdlib.h>

#include "include/util.h"
#include "include/io.h"
#include "include/deserialize.h"

int main(int argc, char **argv) {
    if (argc != 2)
        ERROR_RET_1("Usage: minvm <bytecode>");

    size_t len;
    char *file = read_file(argv[1], &len);

	if (!check_validity(file, len)) {
		free(file);
		ERROR_RET_1("Invalid bytecode file.");
	}

	Chunk out;
	if (!deserialize(file, len, &out)) {
		free(file);
		ERROR_RET_1("Cannot read file.");
	}
	
    free(file);
    return 0;
}

