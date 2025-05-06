#include <stdio.h>
#include <stdlib.h>

#include "util.h"
#include "io.h"
#include "deserialize.h"
#include "vm.h"

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

	Chunk out = {0};
	VM vm = init_vm(&out);
	if (!deserialize(buffer, len, &vm)) {
		free(buffer);
		free_chunk(&out);

		ERROR_RET_1("Cannot read file.");
	}
	
	free_vm(&vm);
	free_chunk(&out);
    free(buffer);
    return 0;
}
