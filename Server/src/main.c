#include "../include/file_handler.h"
#include <time.h>
#include <conio.h>

#define NUM_RECORDS 1000000  // 1 Million

int main() {
    const char* dataFile = "data.bin";
    const char* indexFile = "index.idx";

    // Step 1: Write data to binary file
    writeBinaryFile(dataFile, NUM_RECORDS);

    // Step 2: Create index for fast searching
    createIndexFile(dataFile, indexFile);

    // Step 3: Read first 10 records (sample)
    // readBinaryFile(dataFile);

    // Step 4: Measure search time
    int searchID = 999999;  // Searching the last record
    clock_t start = clock();
    long offset = searchRecord(dataFile, indexFile, searchID);
    clock_t end = clock();

    if (offset != -1) {
        double timeTaken = ((double)(end - start) * 1000.0) / CLOCKS_PER_SEC; // Convert to milliseconds
        printf(" Search Time: %.3f milliseconds\n", timeTaken);
    } else {
        printf(" Record not found!\n");
    }
    getch();
    return 0;
}
