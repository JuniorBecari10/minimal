#ifndef IO_H
#define IO_H

#include "chunk.h"
#include <stdbool.h>

bool read_bytecode(const char *filename, struct chunk *out, struct object **obj_list);

#endif
