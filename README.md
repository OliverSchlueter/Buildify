# Buildify

Easy build and deploy system

## How to build from source

### Windows

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.bat``
3. executable will be located at: ``buildify-core/bin/buildify.exe``

## How to use

**Example: java gradle project**

config.json
```json
{
    "Port": 1337,
    "BuildScriptPath": "build.bat",
    "ArtifactPath": "work/build/libs/MyJar.jar"
}
```

build.bat
````batch
cd work

git fetch
git pull
gradlew shadowJar

cd ../
````

How to start Buildify: ``$ buildify.exe``

### The REST API

``/builds`` - Shows a list of all builds

``/build?id=<number|latest>`` - Shows details about a build

``/download?id=<number>`` - Downloads the output file of a build

``/startBuild`` - Starts the process of creating a new build _(auth required)_

``/deleteBuild?id=<number>`` - Deletes a build _(auth required)_

``/server-stats`` - Shows some statistics about the running server

## Example Web UI

There is an example of how to make a UI for the REST API. Read more [here](buildify-web/README.md).