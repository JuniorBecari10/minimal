#include "chunk.h"

#include <inttypes.h>
#include <stdbool.h>

bool deserialize(const uint8_t *buffer, size_t len, Chunk *out);
