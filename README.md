# Buildify

Easy build and deploy system

## How to build from source

### Windows

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.bat``
3. executable will be at: ``bin/buildify.exe``

### Other

1. figure out yourself, it is a Go project

## How to use

### Run buildify with your config

Example: gradle project

How to start buildify:
``$ buildify.exe -build-script=build.bat -result=work/build/libs/myJar.jar -port=1337``

**build.bat**
````batch
cd work

git fetch
git pull
gradlew shadowJar

cd ../
````

### The API

```/builds``` - Shows all builds

``/build?id=<number|latest>`` - Shows one build

``/startBuild`` - Starts the process of creating a new build

``/download?id=<number>`` - downloads the output file of a build

