#include "chunk.h"

#include <stdlib.h>

void chunk_free(struct chunk *c) {
    free(c->code);
    free(c->constants);
    free(c->metadata);
    
    memset(c, 0, sizeof(*c));
}
