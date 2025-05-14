#ifndef DESERIALIZE_H
#define DESERIALIZE_H

#include "chunk.h"
#include "object.h"

bool deserialize(const char *buffer, size_t buffer_len, struct chunk *out, struct object **obj_list);

#endif
