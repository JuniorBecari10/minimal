#ifndef DESERIALIZE_H
#define DESERIALIZE_H

#include "chunk.h"
#include "object.h"
#include "set.h"

#include <stdbool.h>
#include <stddef.h>

#define TRY(e) if (!e) return false

#define HEADER "MNML"
#define HEADER_LEN 4
#define CHECKSUM_LEN HEADER_LEN

bool deserialize(const char *buffer, size_t buffer_len,
                 struct chunk *out, struct object **obj_list, struct string_set *strings);

#endif
