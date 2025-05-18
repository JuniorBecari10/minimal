#ifndef IO_H
#define IO_H

#include "chunk.h"
#include "object.h"
#include "set.h"

#include <stdbool.h>

bool read_bytecode(const char *file_path,
                   struct chunk *out, struct object **obj_list, struct string_set *strings);

#endif
