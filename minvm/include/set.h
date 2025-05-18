#ifndef SET_H
#define SET_H

#include "string.h"
struct string_set {
    // TODO: make map and make this a wrapper for it.
};

struct string_set string_set_new(void);
struct string *string_set_add(struct string_set *set, struct string str);

#endif
