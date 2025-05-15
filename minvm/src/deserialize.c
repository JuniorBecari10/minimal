#include "deserialize.h"
#include "chunk.h"

#define TRY(e) if (!e) return false

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

}

static bool read_constants(const char *buffer, size_t buffer_len, size_t *counter,
                           struct chunk *out, struct object **obj_list, struct string_set *strings) {

}

static bool read_metadata(const char *buffer, size_t buffer_len, size_t *counter, struct chunk *out) {

}
