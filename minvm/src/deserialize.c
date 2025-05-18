#include "deserialize.h"
#include "chunk.h"
#include "input.h"
#include "codes.h"
#include "object.h"
#include "value.h"

#include <stdlib.h>
#include <string.h>

static bool read_chunk(const char *buffer, size_t buffer_len, size_t *counter,
                       struct chunk *out, struct object **obj_list, struct string_set *strings);

// ---

static bool read_name(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out);

static bool read_code(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out);

static bool read_constants(const char *buffer, size_t buffer_len, size_t *counter,
                           struct chunk *out, struct object **obj_list, struct string_set *strings);

static bool read_metadata(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out);

// ---

static bool read_value(const char *buffer, size_t buffer_len, size_t *counter,
                       struct object **obj_list, struct string_set *strings, struct value *out);

static bool read_meta(const char *buffer, size_t buffer_len, size_t *counter, struct metadata *out);

// ---

bool deserialize(const char *buffer, size_t buffer_len,
                 struct chunk *out, struct object **obj_list, struct string_set *strings) {
    size_t counter = HEADER_LEN; // start after the header
    return read_chunk(buffer, buffer_len, &counter, out, obj_list, strings);
}

static bool read_chunk(const char *buffer, size_t buffer_len, size_t *counter,
                       struct chunk *out, struct object **obj_list, struct string_set *strings) {
    TRY(read_name(buffer, buffer_len, counter, out));
    TRY(read_code(buffer, buffer_len, counter, out));
    TRY(read_constants(buffer, buffer_len, counter, out, obj_list, strings));
    TRY(read_metadata(buffer, buffer_len, counter, out));

    return true;
}
// ---

static bool read_name(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out) {
    uint32_t name_len;
    TRY(read_uint32(buffer, buffer_len, counter, &name_len));

    // see if the buffer can contain all the code
    if (*counter + name_len > buffer_len)
        return false;

    out->name = malloc(name_len);
    
    if (!out->name)
        return false;

    memcpy(out->name, buffer + *counter, name_len);
    *counter += name_len;

    return true;
}

static bool read_code(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out) {
    uint32_t code_len;
    TRY(read_uint32(buffer, buffer_len, counter, &code_len));

    // see if the buffer can contain all the code
    if (*counter + code_len > buffer_len)
        return false;

    out->code = malloc(code_len);
    
    if (!out->code)
        return false;

    memcpy(out->code, buffer + *counter, code_len);
    *counter += code_len;

    return true;
}

static bool read_constants(const char *buffer, size_t buffer_len, size_t *counter,
                           struct chunk *out, struct object **obj_list, struct string_set *strings) {
    uint32_t const_len;
    TRY(read_uint32(buffer, buffer_len, counter, &const_len));

    out->constants = malloc(const_len);
    
    if (!out->constants)
        return false;

    for (struct value *v = out->constants; v < out->constants + const_len; v++) {
        TRY(read_value(buffer, buffer_len, counter, obj_list, strings, v));
    }

    return true;
}

static bool read_metadata(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out) {
    uint32_t metadata_len;
    TRY(read_uint32(buffer, buffer_len, counter, &metadata_len));

    out->metadata = malloc(metadata_len);

    for (size_t i = 0; i < metadata_len; i++) {
        struct metadata metadata;
        TRY(read_meta(buffer, buffer_len, counter, &metadata));

        out->metadata[i] = metadata;
    }

    return true;
}

// ---

static bool read_value(const char *buffer, size_t buffer_len, size_t *counter,
                       struct object **obj_list, struct string_set *strings, struct value *out) {
    uint8_t tag;
    TRY(read_uint8(buffer, buffer_len, counter, &tag));

    switch (tag) {
        case INT_CODE: {
            int32_t num;
            TRY(read_int32(buffer, buffer_len, counter, &num));

            *out = NEW_INT(num);
            return true;
        }

        case FLOAT_CODE: {
            float64 num;
            TRY(read_float64(buffer, buffer_len, counter, &num));

            *out = NEW_FLOAT(num);
            return true;
        }

        case STRING_CODE: {
            struct string str;
            TRY(read_string(buffer, buffer_len, counter, &str));

            struct string *interned = intern_string(str, strings);
            struct obj_string *obj_str = obj_string_new(interned);
            
            TRY(obj_str);

            *out = NEW_OBJECT((struct object *) obj_str);
            return true;
        }

        case BOOL_CODE: {
            uint8_t boolean;
            TRY(read_uint8(buffer, buffer_len, counter, &boolean));

            *out = NEW_BOOL(boolean);
            return true;
        }

        case NIL_CODE: {
            *out = NEW_NIL;
            return true;
        }
        
        case VOID_CODE: {
            *out = NEW_VOID;
            return true;
        }

        case FUNCTION_CODE: {
            uint32_t arity;
            TRY(read_uint32(buffer, buffer_len, counter, &arity));

            uint8_t has_name;
            TRY(read_uint8(buffer, buffer_len, counter, &has_name));

            char *name = NULL;
            if (has_name != 0) {
                struct string str;
                TRY(read_string(buffer, buffer_len, counter, &str));

                // safe, because it is zero-terminated and now 'name' will own the allocation.
                name = str.chars;
            }
            
            struct chunk c;
            TRY(read_chunk(buffer, buffer_len, counter, &c, obj_list, strings));

            struct obj_function *function = obj_function_new(c, arity, name);
            TRY(function);

            struct object *obj = (struct object *) function;
            add_object_to_list(obj, obj_list);

            *out = NEW_OBJECT(obj);
            return true;
        }

        case CLOSURE_CODE: {
            // TODO
            return true;
        }
    }
}

// ---

static bool read_meta(const char *buffer, size_t buffer_len, size_t *counter, struct metadata *out) {
    uint32_t line, col, length;

    TRY(read_uint32(buffer, buffer_len, counter, &line));
    TRY(read_uint32(buffer, buffer_len, counter, &col));
    TRY(read_uint32(buffer, buffer_len, counter, &length));

    *out = (struct metadata) {
        .position = {
            .line = line,
            .col = col,
        },
        .length = length,
    };

    return true;
}
