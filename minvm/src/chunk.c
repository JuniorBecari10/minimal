#include "chunk.h"

#include <stdlib.h>
#include <string.h>

void chunk_free(struct chunk *c) {
    free(c->name);
    free(c->code);
    free(c->constants);
    free(c->metadata);
    
    memset(c, 0, sizeof(*c));
}
