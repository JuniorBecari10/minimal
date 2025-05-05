#include "chunk.h"

#include <stdlib.h>

void free_chunk(Chunk *c) {
    free(c->code);
    
    for (size_t i = 0; i < c->constants.length; i++)
        free_value(&c->constants.data[i]);

    List_Value_free(&c->constants);
    List_Metadata_free(&c->metadata);
}
