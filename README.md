# Buildify

Easy build and deploy system

## How to build from source

### Windows[README.md](README.md)

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.bat``
3. executable will be at: ``buildify-core/bin/buildify.exe``

### Other

1. figure out yourself, it is a Go project

## How to use

### Commandline arguments

``-build-script=<PATH>`` - This script will run everytime a new build is triggered

``-result=<PATH>`` - This is where the downloadable built executable is

``-port=<NUMBER>`` - The API will start on this port

<br>

**Example: java gradle project**

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

### The REST API

```/builds``` - Shows a list of all builds

``/build?id=<number|latest>`` - Shows details about a build

``/download?id=<number>`` - Downloads the output file of a build

``/startBuild`` - Starts the process of creating a new build _(auth required)_

``/deleteBuild?id=<number>`` - Deletes a build _(auth required)_

