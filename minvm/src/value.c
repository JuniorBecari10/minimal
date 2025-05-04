#include "value.h"

void free_value(Value *v) {
    if (IS_OBJ(*v)) {
        // TODO: free the objects
    }

    // primitive values don't need to be freed
}
