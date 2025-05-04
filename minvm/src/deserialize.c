#include "deserialize.h"
#include "value.h"
#include "io.h"

#include <stdio.h>

#define TRY(e) if (!e) return false

static bool read_code(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter);
static bool read_constants(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter);

static bool read_value(const uint8_t *buffer, size_t len, Value *out, size_t *counter);

static bool read_uint8(const uint8_t *buffer, size_t len, size_t *counter, uint8_t *out);
static bool read_uint32(const uint8_t *buffer, size_t len, size_t *counter, uint32_t *out);
static bool read_float64(const uint8_t *buffer, size_t len, size_t *counter, float64 *out);

bool deserialize(const uint8_t *buffer, size_t len, Chunk *out) {
    size_t counter = HEADER_LEN; // to skip the header

    TRY(read_code(buffer, len, out, &counter));
    TRY(read_constants(buffer, len, out, &counter));
    printf("constants ok\n");

    return true;
}

static bool read_code(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter) {
    uint32_t code_len;
    TRY(read_uint32(buffer, len, counter, &code_len));

    if (*counter + code_len > len) return false;
    
    out->code = malloc(code_len);
    if (!out->code) return false;

    memcpy(out->code, buffer + *counter, code_len);
    *counter += code_len;

    return true;
}

static bool read_constants(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter) {
    uint32_t const_len;
    TRY(read_uint32(buffer, len, counter, &const_len));

    out->constants = List_Value_new_with_capacity(const_len);

    for (size_t i = 0; i < const_len; i++) {
        Value value;

        TRY(read_value(buffer, len, &value, counter));
        List_Value_push(&out->constants, value);
    }

    return true;
}

// ---

static bool read_value(const uint8_t *buffer, size_t len, Value *out, size_t *counter) {
    uint8_t tag;
    TRY(read_uint8(buffer, len, counter, &tag));

    switch (tag) {
        case 1: { // Number
            float64 num;
            TRY(read_float64(buffer, len, counter, &num));

            *out = NEW_NUMBER(num);
            return true;
        }

        case 2: { // String

        }

        case 3: { // Bool
            uint8_t boolean;
            TRY(read_uint8(buffer, len, counter, &boolean));

            *out = NEW_BOOL(boolean);
            return true;
        }

        case 4: { // Nil
            *out = NEW_NIL;
            return true;
        }

        case 5: { // Void
            *out = NEW_VOID;
            return true;
        }

        default: {
            fprintf(stderr, "Invalid tag: %u\n", tag);
            return false;
        }
    }
}

// ---

static bool read_uint8(const uint8_t *buffer, size_t len, size_t *counter, uint8_t *out) {
    if (*counter + 1 > len) return false;
    *out = buffer[(*counter)++];

    return true;
}

static bool read_uint32(const uint8_t *buffer, size_t len, size_t *counter, uint32_t *out) {
    if (*counter + 4 > len) return false;

    *out = ((uint32_t)buffer[*counter])           |
           ((uint32_t)buffer[*counter + 1] << 8)  |
           ((uint32_t)buffer[*counter + 2] << 16) |
           ((uint32_t)buffer[*counter + 3] << 24);
    
    *counter += 4;
    return true;
}

static bool read_float64(const uint8_t *buffer, size_t len, size_t *counter, float64 *out) {
    if (*counter + 8 > len) return false;

    uint64_t temp = 
        ((uint64_t) (uint8_t) buffer[(*counter)])           |
        ((uint64_t) (uint8_t) buffer[(*counter) + 1] << 8)  |
        ((uint64_t) (uint8_t) buffer[(*counter) + 2] << 16) |
        ((uint64_t) (uint8_t) buffer[(*counter) + 3] << 24) |
        ((uint64_t) (uint8_t) buffer[(*counter) + 4] << 32) |
        ((uint64_t) (uint8_t) buffer[(*counter) + 5] << 40) |
        ((uint64_t) (uint8_t) buffer[(*counter) + 6] << 48) |
        ((uint64_t) (uint8_t) buffer[(*counter) + 7] << 56);

    memcpy(out, &temp, sizeof(float64));
    *counter += 8;

    return true;
}
