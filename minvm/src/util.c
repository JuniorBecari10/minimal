#include "../include/util.h"

inline static void make_crc32_table(uint32_t *table);

uint32_t compute_checksum(const uint8_t *data, size_t length) {
	uint32_t table[256];
    make_crc32_table(table);

    uint32_t crc = 0xFFFFFFFF;
    for (size_t i = 0; i < length; i++) {
        uint8_t index = (crc ^ data[i]) & 0xFF;
        crc = (crc >> 8) ^ table[index];
    }

    return ~crc;
}

// to make it constant
inline static void make_crc32_table(uint32_t *table) {
    for (uint32_t i = 0; i < 256; i++) {
        uint32_t crc = i;
        for (uint32_t j = 8; j > 0; j--) {
            if (crc & 1)
                crc = (crc >> 1) ^ 0xEDB88320;
            else
                crc >>= 1;
        }
        table[i] = crc;
    }
}
