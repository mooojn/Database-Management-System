#include <stdio.h>
#include <stdlib.h>
#include <windows.h>
#include <conio.h>  // For getch()
#include <time.h> 
#define BUFFER_SIZE 1024
#define THRESHOLD 1024  // Example threshold for mmap usage

void writeToFile(const char* filename, const char* data) {
    size_t dataSize = strlen(data);  

    FILE* file = fopen(filename, "wb");
    if (!file) {
        perror("Error opening file for writing");
        return;
    }

    // Start time
    LARGE_INTEGER start, end, freq;
    QueryPerformanceFrequency(&freq);
    QueryPerformanceCounter(&start);

    fwrite(data, 1, dataSize, file);
    fclose(file);

    // End time
    QueryPerformanceCounter(&end);
    double elapsedTime = (double)(end.QuadPart - start.QuadPart) * 1000.0 / freq.QuadPart;

    printf("Data written successfully.\n");
    printf("Write Time: %.4f ms\n", elapsedTime);
}

void readFromFile(const char* filename, size_t threshold) {
    FILE* file = fopen(filename, "rb");
    if (!file) {
        perror("Error opening file for reading");
        return;
    }

    fseek(file, 0, SEEK_END);
    size_t fileSize = ftell(file);
    rewind(file);

    printf("File size: %zu bytes\n", fileSize);
    printf("Threshold: %zu bytes\n", threshold);

    LARGE_INTEGER start, end, freq;
    QueryPerformanceFrequency(&freq);
    QueryPerformanceCounter(&start);

    if (fileSize > threshold) {
        printf("Using Windows memory-mapped file for reading...\n");

        HANDLE hFile = CreateFileA(filename, GENERIC_READ, FILE_SHARE_READ, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
        if (hFile == INVALID_HANDLE_VALUE) {
            printf("Error opening file for mapping: %lu\n", GetLastError());
            fclose(file);
            return;
        }

        HANDLE hMap = CreateFileMapping(hFile, NULL, PAGE_READONLY, 0, 0, NULL);
        if (!hMap) {
            printf("Error creating file mapping: %lu\n", GetLastError());
            CloseHandle(hFile);
            fclose(file);
            return;
        }

        char* mapped = (char*)MapViewOfFile(hMap, FILE_MAP_READ, 0, 0, fileSize);
        if (!mapped) {
            printf("Mapping failed, falling back to fread: %lu\n", GetLastError());
            CloseHandle(hMap);
            CloseHandle(hFile);
            fclose(file);
            return;
        }

        // printf("Fetched Data:\n");
        // // fwrite(mapped, 1, fileSize, stdout);
        // printf("\n");

        UnmapViewOfFile(mapped);
        CloseHandle(hMap);
        CloseHandle(hFile);
    }
     else {
        printf("Using fread for reading...\n");

        char* buffer = (char*)malloc(fileSize);
        if (!buffer) {
            printf("Memory allocation failed: Requested %zu bytes\n", fileSize);
            fclose(file);
            return;
        }

        fread(buffer, 1, fileSize, file);
        fclose(file);

        // printf("Fetched Data:\n");
        // // fwrite(buffer, 1, fileSize, stdout);
        // printf("\n");

        free(buffer);
    }

    QueryPerformanceCounter(&end);
    double elapsedTime = (double)(end.QuadPart - start.QuadPart) * 1000.0 / freq.QuadPart;
    printf("Read Time: %.4f ms\n", elapsedTime);
}
void searchInFile(const char* filename, const char* keyword) {
    FILE* file = fopen(filename, "rb");
    if (!file) {
        perror("Error opening file for reading");
        return;
    }

    char buffer[BUFFER_SIZE];
    int found = 0;
    long position = 0;

    // Start measuring time
    clock_t start = clock();

    while (fgets(buffer, BUFFER_SIZE, file)) {
        if (strstr(buffer, keyword)) {  // Check if keyword exists in the line
            printf("Match found at position %ld: %s", position, buffer);
            found = 1;
        }
        position = ftell(file);  // Track position in file
    }

    fclose(file);

    // End measuring time
    clock_t end = clock();
    double elapsed_time = ((double)(end - start) / CLOCKS_PER_SEC) * 1000.0; // Convert to milliseconds

    if (!found) {
        printf("No matches found for '%s'\n", keyword);
    }

    printf("Search Time: %.4f ms\n", elapsed_time);
}