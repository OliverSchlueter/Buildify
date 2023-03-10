package main

import (
	"Buildify/builds"
	"flag"
	"log"
)

var (
	buildScriptPath *string
	resultPath      *string
)

func main() {
	buildScriptPath = flag.String("build-script", "build.bat", "path to the build script")
	resultPath = flag.String("result", "work/build/libs/StackPP.jar", "path to the result executable file")
	flag.Parse()

	err := builds.LoadBuildsFile("builds/")
	if err != nil {
		log.Println("Could not load build metadata file")
	}

	builds.BuildBuild(buildScriptPath, resultPath)

	err = builds.SaveBuildsFile("builds/")
	if err != nil {
		log.Println("Could not save build metadata")
	}
}
