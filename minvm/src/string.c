#include "string.h"

#include <string.h>
#include <stdlib.h>

void string_free(String *s) {
    free(s->chars);
    memset(s, 0, sizeof(*s));
}
