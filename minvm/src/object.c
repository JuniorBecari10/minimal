#include "value.h"
#include "object.h"

#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

bool is_obj_type(Value value, ObjType type) {
    return IS_OBJ(value) && AS_OBJECT(value)->type == type;
}

void object_free(Object *obj) {
    // memset should be according to the variant to fill the entire struct.
    // obviously, this switch must be exhaustive, otherwise a memory leak may happen.
	switch (obj->type) {
        case OBJ_STRING: {
            ObjString *s = (ObjString *) obj;

            free(s->str);
            memset(obj, 0, sizeof(*s));
        }

        case OBJ_FUNCTION: {
            ObjFunction *fn = (ObjFunction *) obj;

            chunk_free(&fn->chunk);
            free(&fn->name);
            
            memset(obj, 0, sizeof(*fn));
        }
    }
}
