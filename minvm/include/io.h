#include <stdlib.h>
#include <stdbool.h>

// returns NULL if error.
char *read_file(const char *path, size_t *output_len);
bool check_validity(char *file, size_t len);
