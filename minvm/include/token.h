#ifndef TOKEN_H
#define TOKEN_H

#include <stdint.h>

typedef struct {
	uint32_t line;
	uint32_t col;
} Position;

typedef struct {
	Position position;
	uint32_t length;
} Metadata;

#endif
