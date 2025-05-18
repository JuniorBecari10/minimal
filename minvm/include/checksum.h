#ifndef CHECKSUM_H
#define CHECKSUM_H

#include <stddef.h>
#include <inttypes.h>

uint32_t compute_checksum(const char *data, size_t length);

#endif
