#include "../include/file_handler.h"

void writeBinaryFile(const char* filename, int numRecords) {
    FILE* file = fopen(filename, "wb");
    if (!file) {
        perror("Error opening file for writing");
        return;
    }

    for (int i = 0; i < numRecords; i++) {
        Record rec = { i, "TestData", (float)i * 1.5 };
        fwrite(&rec, sizeof(Record), 1, file);
    }

    fclose(file);
    printf("✅ %d records written to %s\n", numRecords, filename);
}

void readBinaryFile(const char* filename) {
    FILE* file = fopen(filename, "rb");
    if (!file) {
        perror("Error opening file for reading");
        return;
    }

    Record rec;
    // printf("\n🔍 Reading records from %s:\n", filename);
    // while (fread(&rec, sizeof(Record), 1, file)) {
    //     printf("ID: %d, Name: %s, Value: %.2f\n", rec.id, rec.name, rec.value);
    // }

    fclose(file);
}

void createIndexFile(const char* dataFile, const char* indexFile) {
    FILE* data = fopen(dataFile, "rb");
    FILE* index = fopen(indexFile, "wb");

    if (!data || !index) {
        perror("Error opening files");
        return;
    }

    Record rec;
    long offset = 0;

    while (fread(&rec, sizeof(Record), 1, data)) {
        fwrite(&offset, sizeof(long), 1, index);
        offset = ftell(data);
    }

    fclose(data);
    fclose(index);
    printf(" Index file %s created.\n", indexFile);
}

long searchRecord(const char* dataFile, const char* indexFile, int recordID) {
    FILE* index = fopen(indexFile, "rb");
    if (!index) {
        perror("Error opening index file");
        return -1;
    }

    long offset;
    fseek(index, recordID * sizeof(long), SEEK_SET);
    fread(&offset, sizeof(long), 1, index);
    fclose(index);

    FILE* data = fopen(dataFile, "rb");
    if (!data) {
        perror("Error opening data file");
        return -1;
    }

    fseek(data, offset, SEEK_SET);
    Record rec;
    fread(&rec, sizeof(Record), 1, data);
    fclose(data);

    printf("\n Found Record - ID: %d, Name: %s, Value: %.2f\n", rec.id, rec.name, rec.value);
    return offset;
}
// long searchRecord(const char* dataFile, const char* indexFile, int recordID) {
//     FILE* data = fopen(dataFile, "rb");
//     if (!data) {
//         perror("Error opening data file");
//         return -1;
//     }

//     Record rec;
//     long offset = -1;

//     clock_t start = clock();  // Start timing

//     while (fread(&rec, sizeof(Record), 1, data)) {
//         if (rec.id == recordID) {
//             offset = ftell(data) - sizeof(Record);
//             printf("\n✅ Found Record - ID: %d, Name: %s, Value: %.2f\n", rec.id, rec.name, rec.value);
//             break;
//         }
//     }

//     fclose(data);

//     clock_t end = clock();  // End timing
//     double timeTaken = ((double)(end - start) * 1000.0) / CLOCKS_PER_SEC;
//     printf("⏳ Search Time: %.3f ms\n", timeTaken);

//     if (offset == -1) {
//         printf("❌ Record with ID %d not found.\n", recordID);
//     }

//     return offset;
// }

