@ECHO OFF
CLS

:: go into sourcecode directory
CD buildify-core

:: run tests
go test .

:: build the program
go build -buildmode exe -o bin/buildify.exe

:: go back to main directory
CD ../