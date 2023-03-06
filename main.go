package main

import (
	"flag"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"os/exec"
)

func main() {
	buildScriptPath := flag.String("build-script", "build.bat", "path to the build script")
	resultPath := flag.String("result", "work/build/libs/StackPP.jar", "path to the result executable file")
	flag.Parse()

	// create working directory
	err := createWorkingDir()
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	dir := os.Getenv("=D:") + "\\"

	// build project
	err, built := buildProject(dir + *buildScriptPath)
	if err != nil {
		log.Fatal(err)
	}

	if !built {
		log.Fatal("Could not build project")
	}

	// get result file
	err, resultFile := getResultFile(*resultPath)
	if err != nil {
		log.Fatal(err)
	}

	defer resultFile.Close()

	println("Result file: " + resultFile.Name())

	// get git hash
	hash := getGitHash()
	println("Hash: " + hash)
}

func createWorkingDir() error {
	log.Println("Creating working directory (./work/)")
	err := os.Mkdir("work", os.ModePerm)
	return err
}

func buildProject(buildScript string) (error, bool) {
	log.Println("Building the project (" + buildScript + ")")
	command := exec.Command(buildScript)
	//command.Stdout = os.Stdout

	err := command.Run()
	if err != nil {
		return err, false
	}

	success := command.ProcessState.ExitCode() == 0
	return nil, success
}

func getResultFile(path string) (error, *os.File) {
	log.Println("Getting the result file")
	file, err := os.OpenFile(path, 0, os.ModePerm)
	if err != nil {
		log.Fatal("Could not find result file: " + path)
	}

	return err, file
}

func getGitHash() string {
	log.Println("Getting the commit hash")
	gitApp, err := git.PlainOpen("work")
	if err != nil {
		log.Fatal(err)
	}

	gitHead, err := gitApp.Head()
	if err != nil {
		log.Fatal(err)
	}

	hash := gitHead.Hash()

	return hash.String()
}
