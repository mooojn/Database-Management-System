#ifndef FILE_HANDLER_H
#define FILE_HANDLER_H

#include <stdio.h>
#include <stdlib.h>

#include <fcntl.h>

#include <time.h>

typedef struct {
    int id;
    char name[20];
    float value;
} Record;

void writeBinaryFile(const char* filename, int numRecords);
void readBinaryFile(const char* filename);
void createIndexFile(const char* dataFile, const char* indexFile);
long searchRecord(const char* dataFile, const char* indexFile, int recordID);

#endif
