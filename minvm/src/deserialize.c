#include "deserialize.h"
#include "chunk.h"

#define TRY(e) if (!e) return false

static bool read_chunk(const char *buffer, size_t buffer_len, struct chunk *out, struct object **obj_list, size_t *counter);

bool deserialize(const char *buffer, size_t buffer_len, struct chunk *out, struct object **obj_list) {

}
