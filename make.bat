@ECHO OFF

CLS
:: run tests
go test .

:: build the program
go build -o bin/buildify.exe .