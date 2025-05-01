#include "lists.h"
#include <stdint.h>

typedef struct {
	char *code;
	List_Value constants;
	List_Metadata metadata;
} Chunk;

