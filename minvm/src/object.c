#include "value.h"
#include "object.h"

#include <stdlib.h>
#include <stdbool.h>

bool is_obj_type(Value value, ObjType type) {
    return IS_OBJ(value) && AS_OBJECT(value)->type == type;
}

void free_object(Object *obj) {
	switch (obj->type) {
        case OBJ_STRING: {
            ObjString *s = (ObjString *) obj;
            free(s->str);
        }
    }
}
