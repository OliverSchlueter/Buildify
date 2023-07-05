@ECHO OFF
CLS

:: go into sourcecode directory
CD src

:: run tests
go test .

:: build the program
go build -buildmode exe -o ../bin/buildify.exe

:: go back to main directory
CD ../

XCOPY "static" "bin/static/" /E /I /Y