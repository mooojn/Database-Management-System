#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

typedef struct {
    int id;
    float score;
} Record;

#define TOTAL_RECORDS 1000000  // 1 million records

// Function to search in binary file
void searchBinaryFile(int search_id) {
    FILE *file = fopen("data.bin", "rb");
    if (!file) {
        printf("Error opening binary file!\n");
        return;
    }

    Record r;
    clock_t start = clock();
    
    // Jump directly to the record (assuming sequential IDs)
    fseek(file, search_id * sizeof(Record), SEEK_SET);
    fread(&r, sizeof(Record), 1, file);

    clock_t end = clock();
    fclose(file);

    if (r.id == search_id)
        printf("Binary Search: ID=%d, Score=%.2f (Time: %lf sec)\n", r.id, r.score, (double)(end - start) / CLOCKS_PER_SEC);
    else
        printf("Binary Search: ID not found!\n");
}

// Function to search in CSV file
void searchCSVFile(int search_id) {
    FILE *file = fopen("data.csv", "r");
    if (!file) {
        printf("Error opening CSV file!\n");
        return;
    }

    char line[100];
    int id;
    float score;
    clock_t start = clock();
    
    // Linear search in CSV file
    while (fgets(line, sizeof(line), file)) {
        sscanf(line, "%d,%f", &id, &score);
        if (id == search_id) {
            clock_t end = clock();
            printf("CSV Search: ID=%d, Score=%.2f (Time: %lf sec)\n", id, score, (double)(end - start) / CLOCKS_PER_SEC);
            fclose(file);
            return;
        }
    }
    
    fclose(file);
    printf("CSV Search: ID not found!\n");
}

int main() {
    int search_id = 900000;  // Change this to test different IDs

    searchBinaryFile(search_id);
    searchCSVFile(search_id);

    return 0;
}
