# Buildify

Easy build and deploy system

## How to build from source

### Windows

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.bat``
3. executable will be located at: ``buildify-core/bin/buildify.exe``

### Linux

1. clone the repo ``git clone https://github.com/OliverSchlueter/Buildify.git``
2. run ``make.sh``
3. executable will be located at: ``buildify-core/bin/buildify``


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

``/api/builds`` - Shows a list of all builds

``/api/build/<id|latest>`` - Shows details about a build

``/api/download/<id>`` - Downloads the output file of a build

``/api/startBuild`` - Starts the process of creating a new build _(auth required)_

``/api/deleteBuild/<id>`` - Deletes a build _(auth required)_

``/api/server-stats`` - Shows some statistics about the running server

## Example Web UI

There is an example of how to make a UI for the REST API. Read more [here](buildify-web/README.md).