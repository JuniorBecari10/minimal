#ifndef INPUT_H
#define INPUT_H

#include "value.h"

#include <stdbool.h>
#include <stddef.h>
#include <inttypes.h>

bool read_uint8(const char *buffer, size_t buffer_len, size_t *counter, uint8_t *out);
bool read_uint32(const char *buffer, size_t buffer_len, size_t *counter, uint32_t *out);
bool read_float64(const char *buffer, size_t buffer_len, size_t *counter, float64 *out);

#endif
