package builds

import (
	"Buildify/util"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func BuildBuild(buildScriptPath, resultPath *string) Build {
	// create working directory
	err := createWorkingDir()
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	dir := os.Getenv("=D:") + "\\"

	// build project
	buildId := len(Builds) + 1
	buildTime := time.Now().UnixMilli()

	err, built := buildProject(dir + *buildScriptPath)
	if err != nil {
		log.Fatal(err)
	}

	if !built {
		log.Fatal("Could not build project")
	}

	// get result file
	err, resultFile := getResultFile(*resultPath, buildId)
	if err != nil {
		log.Fatal(err)
	}

	defer resultFile.Close()

	// get git hash
	gitHash, gitMessage := getGitInfo()

	println("----------------------------------------------------------")
	println("Id: " + strconv.Itoa(buildId))
	println("Time: " + time.UnixMilli(buildTime).String())
	println("File: " + resultFile.Name())
	println("Hash: " + gitHash)
	println("Message: " + gitMessage)
	println("----------------------------------------------------------")

	b := Create(buildId, buildTime, gitHash, gitMessage, resultFile.Name())

	return b
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

func getResultFile(path string, buildId int) (error, *os.File) {
	log.Println("Getting the result file")
	file, err := os.OpenFile(path, 0, os.ModePerm)
	if err != nil {
		log.Fatal("Could not find result file: " + path)
	}

	defer file.Close()

	stat, _ := file.Stat()
	fileName := stat.Name()

	// copy file
	buildDir := "builds/"
	err = os.Mkdir(buildDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err, nil
	}

	buildDir = buildDir + "build-" + strconv.Itoa(buildId) + "/"
	err = os.Mkdir(buildDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err, nil
	}

	copiedFile, err := os.Create(buildDir + fileName)
	if err != nil {
		log.Fatal("Could not copy result file1")
	}

	err = util.FastCopyFile(file, copiedFile)
	if err != nil {
		log.Fatal("Could not copy result file")
	}

	return err, copiedFile
}

func getGitInfo() (string, string) {
	log.Println("Getting information about the latest commit")
	gitApp, err := git.PlainOpen("work")
	if err != nil {
		log.Fatal(err)
	}

	head, err := gitApp.Head()
	if err != nil {
		log.Fatal(err)
	}

	hash := head.Hash()

	commit, err := gitApp.CommitObject(hash)
	if err != nil {
		log.Fatal(err)
	}

	message := commit.Message
	message = message[0 : len(message)-1] // remove last '\n'

	return hash.String(), message
}
