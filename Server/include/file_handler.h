#ifndef FILE_HANDLER_H
#define FILE_HANDLER_H

#include <stdio.h>
#include <stdlib.h>

#include <fcntl.h>

#include <time.h>

void writeBinaryFile(int n) ;
void writeCSVFile(int n) ;
void searchBinaryFile(int search_id) ;
void searchCSVFile(int search_id) ;
void readF();
void writeF();

#endif