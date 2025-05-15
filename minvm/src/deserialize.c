#include "deserialize.h"
#include "chunk.h"
#include "input.h"

#include <stdlib.h>
#include <string.h>

static bool read_chunk(const char *buffer, size_t buffer_len, size_t *counter,
                       struct chunk *out, struct object **obj_list, struct string_set *strings);

// ---

static bool read_code(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out);

static bool read_constants(const char *buffer, size_t buffer_len, size_t *counter,
                           struct chunk *out, struct object **obj_list, struct string_set *strings);

static bool read_metadata(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out);

// ---

bool deserialize(const char *buffer, size_t buffer_len,
                 struct chunk *out, struct object **obj_list, struct string_set *strings) {
    size_t counter = HEADER_LEN;
    return read_chunk(buffer, buffer_len, &counter, out, obj_list, strings);
}

static bool read_chunk(const char *buffer, size_t buffer_len, size_t *counter,
                       struct chunk *out, struct object **obj_list, struct string_set *strings) {
    TRY(read_code(buffer, buffer_len, counter, out));
    TRY(read_constants(buffer, buffer_len, counter, out, obj_list, strings));
    TRY(read_metadata(buffer, buffer_len, counter, out));

    return true;
}
// ---

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

}

static bool read_metadata(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out) {

}
