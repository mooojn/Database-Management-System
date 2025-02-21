#include <stdio.h>
#include <time.h>

typedef struct {
    int id;
    float score;
} Record;

void writeBinaryFile(int n) {
    FILE *file = fopen("data.bin", "wb");
    Record r;
    for (int i = 0; i < n; i++) {
        r.id = i;
        r.score = i * 1.5;
        fwrite(&r, sizeof(Record), 1, file);
    }
    fclose(file);
}

void writeCSVFile(int n) {
    FILE *file = fopen("data.csv", "w");
    for (int i = 0; i < n; i++) {
        fprintf(file, "%d,%.2f\n", i, i * 1.5);
    }
    fclose(file);
}

int main() {
    int n = 1000000;
    clock_t start, end;

    start = clock();
    writeBinaryFile(n);
    end = clock();
    printf("Binary Write Time: %lf seconds\n", (double)(end - start) / CLOCKS_PER_SEC);

    start = clock();
    writeCSVFile(n);
    end = clock();
    printf("CSV Write Time: %lf seconds\n", (double)(end - start) / CLOCKS_PER_SEC);

    return 0;
}
