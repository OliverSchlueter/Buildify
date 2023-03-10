@ECHO OFF

CLS
:: run tests
go test .

:: build the program
go build -buildmode exe -o bin/buildify.exe