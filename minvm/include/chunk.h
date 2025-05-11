#ifndef CHUNK_H
#define CHUNK_H

#include "lists.h"
#include <stdint.h>

typedef struct {
	uint8_t *code;
	List_Value constants;
	List_Metadata metadata;
} Chunk;

void chunk_free(Chunk *c);

#endif
