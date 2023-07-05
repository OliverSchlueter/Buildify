clear

cd src

go build -buildmode exe -o ../bin/buildify

cd ..

cp -r "./static/" "./bin/static/"