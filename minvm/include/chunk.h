#include "lists.h"
#include <stdint.h>

typedef struct {
	uint8_t *code;
	List_Value constants;
	List_Metadata metadata;
} Chunk;
