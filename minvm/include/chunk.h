#ifndef CHUNK_H
#define CHUNK_H

#include "extra.h"

#include <stddef.h>
#include <inttypes.h>

struct chunk {
    uint8_t *code;

    struct value *constants;
    size_t constants_len;

    struct metadata *metadata;
    size_t metadata_len;
};

void chunk_free(struct chunk *c);

#endif
