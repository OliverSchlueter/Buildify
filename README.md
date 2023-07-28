# Buildify

Lightweight build and deployment system.<br>
_Made for educational purposes, do not use in production._

## How to build from source

### Windows

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.bat``
3. executable will be located at: ``bin/buildify.exe``

### Linux

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.sh``
3. executable will be located at: ``bin/buildify``


## How to use

1. download buildify.exe from github releases or build yourself
2. run buildify.exe and let it generate some files
3. stop buildify
4. edit config.json as you wish
5. create your build script (and set up the work directory)
6. start buildify

**Example config and build script for a java gradle project**

config.json
```json
{
    "Port": 1337,
    "BuildScriptPath": "build.bat",
    "ArtifactPath": "work/build/libs/MyJar.jar"
}
```

build.bat
```batch
cd work

git fetch
git pull

gradlew build

cd ../
```

### The REST API

``/api/builds`` - Shows a list of all builds

``/api/build/<id|latest>`` - Shows details about a build

``/api/download/<id>`` - Downloads the output file of a build

``/api/startBuild`` - Starts the process of creating a new build _(auth required)_

``/api/deleteBuild/<id>`` - Deletes a build _(auth required)_

``/api/server-stats`` - Shows some statistics about the running server

``/api/build-script`` (GET) - Shows the current build script

``/api/build-script`` (POST) - Sets the current build script (put the new script in the body)

``/api/artifact-file-path`` (GET) - Shows the current artifact file path

``/api/artifact-file-path`` (POST) - Sets the current artifact file path (put the new path in the body)