#include "deserialize.h"
#include "object.h"
#include "string.h"
#include "value.h"
#include "io.h"
#include "vm.h"

#include <stdio.h>
#include <string.h>

#define TRY(e) if (!e) return false

static bool read_code(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter);
static bool read_constants(const uint8_t *buffer, size_t len, VM *vm, size_t *counter);
static bool read_metadata(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter);

static bool read_value(const uint8_t *buffer, size_t len, Value *out, VM *vm, size_t *counter);
static bool read_meta(const uint8_t *buffer, size_t len, Metadata *out, size_t *counter);

static bool read_uint8(const uint8_t *buffer, size_t len, size_t *counter, uint8_t *out);
static bool read_uint32(const uint8_t *buffer, size_t len, size_t *counter, uint32_t *out);
static bool read_float64(const uint8_t *buffer, size_t len, size_t *counter, float64 *out);
static bool read_string(const uint8_t *buffer, size_t len, size_t *counter, String *out);

bool deserialize(const uint8_t *buffer, size_t len, VM *vm) {
    size_t counter = HEADER_LEN; // to skip the header

    TRY(read_code(buffer, len, vm->chunk, &counter));
    TRY(read_constants(buffer, len, vm, &counter));
    TRY(read_metadata(buffer, len, vm->chunk, &counter));

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

static bool read_constants(const uint8_t *buffer, size_t len, VM *vm, size_t *counter) {
    uint32_t const_len;
    TRY(read_uint32(buffer, len, counter, &const_len));

    vm->chunk->constants = List_Value_new_with_capacity(const_len);

    for (size_t i = 0; i < const_len; i++) {
        Value value;

        TRY(read_value(buffer, len, &value, vm, counter));
        List_Value_push(&vm->chunk->constants, value);
    }

    return true;
}

static bool read_metadata(const uint8_t *buffer, size_t len, Chunk *out, size_t *counter) {
    uint32_t metadata_len;
    TRY(read_uint32(buffer, len, counter, &metadata_len));

    out->metadata = List_Metadata_new_with_capacity(metadata_len);

    for (size_t i = 0; i < metadata_len; i++) {
        Metadata metadata;

        TRY(read_meta(buffer, len, &metadata, counter));
        List_Metadata_push(&out->metadata, metadata);
    }

    return true;
}

// ---

static bool read_value(const uint8_t *buffer, size_t len, Value *out, VM *vm, size_t *counter) {
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
            String string;
            TRY(read_string(buffer, len, counter, &string));
            
            String *ptr = intern_string(vm, string);
            ObjString *str = allocate_object(vm, sizeof(ObjString), OBJ_STRING);
            str->str = ptr;

            *out = NEW_OBJECT(str);
            return true;
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

static bool read_meta(const uint8_t *buffer, size_t len, Metadata *out, size_t *counter) {
    uint32_t line, col, length;

    TRY(read_uint32(buffer, len, counter, &line));
    TRY(read_uint32(buffer, len, counter, &col));
    TRY(read_uint32(buffer, len, counter, &length));

    *out = (Metadata) {
        .position = {
            .line = line,
            .col = col,
        },
        .length = length,
    };

    return true;
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

static bool read_string(const uint8_t *buffer, size_t len, size_t *counter, String *out) {
    uint32_t string_len;
    TRY(read_uint32(buffer, len, counter, &string_len));

    // Check for buffer overflow
    if (*counter + string_len > len)
        return false;

    // String + null terminator
    char *s = malloc(string_len + 1);
    if (s == NULL)
        return false;

    memcpy(s, buffer + *counter, string_len);
    s[string_len] = '\0';

    *counter += string_len;
    *out = (String) {
        .chars = s,
        .length = string_len,
    };

    return true;
}
