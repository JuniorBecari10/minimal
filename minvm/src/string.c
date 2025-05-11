#include "string.h"
#include <stdlib.h>

void string_free(String *s) {
    free(s->chars);
}
