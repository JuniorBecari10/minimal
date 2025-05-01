#include <stdint.h>

typedef struct {
	uint32_t line;
	uint32_t col;
} Position;

typedef struct {
	char *code;
	List_Value constants;
	List_Metadata metadata;
} Chunk;

typedef struct {
	Position position;
	uint32_t length;
} Metadata;

