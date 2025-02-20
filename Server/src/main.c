#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <conio.h>
#include "../include/file_handler.h"

#define THRESHOLD 1024 
#define NUM_ENTRIES 100000

int main() {
    const char* filename = "data.bin";

    // Allocate memory for storing data
    size_t bufferSize = NUM_ENTRIES * 20;  // Approximate size
    char* finalData = (char*)malloc(bufferSize);
    if (!finalData) {
        printf("Memory allocation failed.\n");
        return 1;
    }
    finalData[0] = '\0';  // Initialize empty string

    // Populate the buffer with data
    for (int i = 0; i < NUM_ENTRIES; i++) {
        char entry[32];  // Enough space for "Entry Number: X\n"
        sprintf(entry, "Entry Number: %d\n", i);
        strcat(finalData, entry);
    }
    
    // Write data to file
    writeToFile(filename, finalData);

    // Read data from file
    readFromFile(filename, THRESHOLD);

    // Free allocated memory
    free(finalData);

    const char* searchTerm = "Entry Number: 99999"; 

    searchInFile(filename, searchTerm);
    printf("\nPress any key to exit...");
    getch();  // Pause screen
    return 0;
}
