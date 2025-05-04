#include "chunk.h"

#include <inttypes.h>
#include <stdbool.h>

bool deserialize(char *file, size_t len, Chunk *out);
