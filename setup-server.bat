@echo off
mkdir server\bin 2>nul

:: Compile file_handler and main
gcc -o server\bin\main.exe server\src\main.c server\src\file_handler.c -Iserver\src\core

:: Compile server with mongoose
gcc -o server\bin\server.exe server\src\server.c server\src\mongoose.c -lws2_32 -D_CRT_RAND_S -Iserver\src

echo Build completed!
pause
