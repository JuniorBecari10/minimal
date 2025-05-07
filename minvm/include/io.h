#ifndef IO_H
#define IO_H

#include <stdlib.h>
#include <inttypes.h>
#include <stdbool.h>

#define HEADER "MNML"
#define HEADER_LEN 4
#define CHECKSUM_LEN HEADER_LEN

// returns NULL if error.
uint8_t *read_file(const char *path, size_t *output_len);
bool check_validity(const uint8_t *buffer, size_t len);

#endif
