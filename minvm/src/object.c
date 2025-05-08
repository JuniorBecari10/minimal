#include "value.h"
#include "object.h"

#include <stdbool.h>

bool is_obj_type(Value value, ObjType type) {
    return IS_OBJ(value) && AS_OBJECT(value)->type == type;
}

void free_object(Object *obj) {
	// TODO: free object
}
