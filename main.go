package main

import (
	"Buildify/builds"
	"flag"
)

var (
	buildScriptPath *string
	resultPath      *string
)

func main() {
	buildScriptPath = flag.String("build-script", "build.bat", "path to the build script")
	resultPath = flag.String("result", "work/build/libs/StackPP.jar", "path to the result executable file")
	flag.Parse()

	builds.BuildBuild(buildScriptPath, resultPath)
}
