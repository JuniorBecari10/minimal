#include "input.h"
#include "deserialize.h"
#include "string.h"

#include <stdlib.h>
#include <string.h>

bool read_uint8(const char *buffer, size_t buffer_len, size_t *counter, uint8_t *out) {
    if (*counter + 1 > buffer_len) return false;
        *out = buffer[(*counter)++];

        return true;
}

bool read_int32(const char *buffer, size_t buffer_len, size_t *counter, int32_t *out) {
    if (*counter + 4 > buffer_len)
        return false;

    *out = ((int32_t) buffer[*counter])           |
           ((int32_t) buffer[*counter + 1] << 8)  |
           ((int32_t) buffer[*counter + 2] << 16) |
           ((int32_t) buffer[*counter + 3] << 24);
    
    *counter += 4;
    return true;
}

bool read_uint32(const char *buffer, size_t buffer_len, size_t *counter, uint32_t *out) {
    if (*counter + 4 > buffer_len)
        return false;

    *out = ((uint32_t) buffer[*counter])           |
           ((uint32_t) buffer[*counter + 1] << 8)  |
           ((uint32_t) buffer[*counter + 2] << 16) |
           ((uint32_t) buffer[*counter + 3] << 24);
    
    *counter += 4;
    return true;
}

bool read_float64(const char *buffer, size_t buffer_len, size_t *counter, float64 *out) {
    if (*counter + 8 > buffer_len)
        return false;

    uint64_t temp =
        ((uint64_t) buffer[(*counter)])           |
        ((uint64_t) buffer[(*counter) + 1] << 8)  |
        ((uint64_t) buffer[(*counter) + 2] << 16) |
        ((uint64_t) buffer[(*counter) + 3] << 24) |
        ((uint64_t) buffer[(*counter) + 4] << 32) |
        ((uint64_t) buffer[(*counter) + 5] << 40) |
        ((uint64_t) buffer[(*counter) + 6] << 48) |
        ((uint64_t) buffer[(*counter) + 7] << 56);

    memcpy(out, &temp, sizeof(float64));
    *counter += 8;

    return true;
}

bool read_string(const char *buffer, size_t buffer_len, size_t *counter, struct string *out) {
    uint32_t string_len;
    TRY(read_uint32(buffer, buffer_len, counter, &string_len));

    if (*counter + string_len > buffer_len)
        return false;

    *out = string_new(buffer + *counter, string_len);
    return true;
}
