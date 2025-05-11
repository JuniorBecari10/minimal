#include "value.h"
#include "object.h"

void value_free(Value *v) {
    if (IS_OBJ(*v))
		object_free(AS_OBJECT(*v));

    // primitive values don't need to be freed
}
