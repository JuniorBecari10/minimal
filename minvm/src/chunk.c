#include "chunk.h"

#include <stdlib.h>

void free_chunk(Chunk *c) {
    free(c->code);
    
    for (size_t i = 0; i < c->constants.length; i++)
        free_value(&c->constants.data[i]);

    // don't free constants, because they are objects in the VM, and it frees them
    List_Metadata_free(&c->metadata);
}
