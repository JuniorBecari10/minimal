#include "chunk.h"

#include <stdlib.h>

void chunk_free(Chunk *c) {
    free(c->code);
    
    for (size_t i = 0; i < c->constants.length; i++)
        value_free(&c->constants.data[i]);

    // don't free constants, because they are objects in the VM, and it frees them
    List_Metadata_free(&c->metadata);
    memset(c, 0, sizeof(*c));
}
