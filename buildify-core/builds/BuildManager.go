package builds

import (
	"Buildify/util"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
)

var IsBuilding bool = false

func CreateBuild(buildScriptPath, resultPath *string) (error, *Build) {
	// create working directory
	err := createWorkingDir()
	if err != nil && !os.IsExist(err) {
		log.Println(err)
		return err, nil
	}

	dir := os.Getenv("=D:") + "\\"

	// build project
	var buildId int
	if len(Builds) == 0 {
		buildId = 1
	} else {
		buildId = Builds[len(Builds)-1].Id + 1
	}

	buildTime := time.Now().UnixMilli()

	err, built := buildProject(dir + *buildScriptPath)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	if !built {
		return errors.New("could not build project"), nil
	}

	// get result file
	err, resultFile := getResultFile(*resultPath, buildId)
	if err != nil {
		log.Println(err)
		return err, nil
	}

	defer resultFile.Close()

	// get git hash
	err, gitHash, gitMessage := getGitInfo()
	if err != nil {
		log.Println(err)
		return err, nil
	}

	log.Println("Finished build (#" + strconv.Itoa(buildId) + ")")
	b := Create(buildId, buildTime, gitHash, gitMessage, 0, resultFile.Name())

	return nil, b
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
		return errors.New("Could not find result file: " + path), nil
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
		return errors.New("could not copy result file"), nil
	}

	err = util.FastCopyFile(file, copiedFile)
	if err != nil {
		return errors.New("could not copy result file"), nil
	}

	return err, copiedFile
}

func getGitInfo() (error, string, string) {
	log.Println("Getting information about the latest commit")
	gitApp, err := git.PlainOpen("work")
	if err != nil {
		log.Println(err)
		return err, "", ""
	}

	head, err := gitApp.Head()
	if err != nil {
		log.Println(err)
		return err, "", ""
	}

	hash := head.Hash()

	commit, err := gitApp.CommitObject(hash)
	if err != nil {
		log.Println(err)
		return err, "", ""
	}

	message := commit.Message
	message = message[0 : len(message)-1] // remove last '\n'

	return nil, hash.String(), message
}
