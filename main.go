package main

import (
	"Buildify/api"
	"Buildify/builds"
	"Buildify/util"
	"flag"
	"log"
)

func main() {
	util.BuildScriptPath = flag.String("build-script", "build.bat", "path to the build script")
	util.ResultPath = flag.String("result", "work/build/libs/StackPP.jar", "path to the result executable file")
	port := flag.Int("port", 1337, "port for the REST API")
	flag.Parse()

	err := builds.LoadBuildsFile("builds/")
	if err != nil {
		log.Println("Could not load build metadata file")
	}

	//builds.BuildBuild(buildScriptPath, resultPath)

	api.Start(*port)
}
