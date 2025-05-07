#include "../include/value.h"
#include "../include/object.h"

void free_value(Value *v) {
    if (IS_OBJ(*v))
		free_object(AS_OBJECT(*v));

    // primitive values don't need to be freed
}
