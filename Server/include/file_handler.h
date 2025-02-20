#ifndef FILE_HANDLER_H
#define FILE_HANDLER_H

#include <stdio.h>
#include <stdlib.h>

void writeToFile(const char* filename, const char* data);
char* readFromFile(const char* filename, size_t threshold);
void searchInFile(const char* filename, const char* keyword);
#endif
