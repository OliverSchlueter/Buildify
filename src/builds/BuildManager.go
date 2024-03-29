package builds

import (
	"Buildify/config"
	"Buildify/util"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
)

var IsBuilding bool = false

func CreateBuild() (error, *Build) {
	util.LogInfo("------------------------------------------------------")
	util.LogInfo("")
	util.LogInfo("Creating new build")
	util.LogInfo("")

	startTime := time.Now()

	// create working directory
	err := createWorkingDir()
	if err != nil && !os.IsExist(err) {
		util.LogInfo(err.Error())
		return err, nil
	}

	// build project
	var buildId int
	if len(Builds) == 0 {
		buildId = 1
	} else {
		buildId = Builds[len(Builds)-1].Id + 1
	}

	dir := os.Getenv("=D:") + "\\"
	err, built := buildProject(dir + config.CurrentConfig.BuildScriptPath)
	if err != nil {
		util.LogInfo(err.Error())
		return err, nil
	}

	if !built {
		return errors.New("could not build project"), nil
	}

	// get artifact file
	artifactFile, err := getArtifactFile(config.CurrentConfig.ArtifactPath, buildId)
	if err != nil {
		util.LogInfo(err.Error())
		return err, nil
	}

	defer artifactFile.Close()

	// get git hash
	err, gitHash, gitMessage := getGitInfo()
	if err != nil {
		util.LogInfo(err.Error())
		return err, nil
	}

	buildTime := time.Now().UnixMilli()

	b := Create(buildId, buildTime, gitHash, gitMessage, 0, artifactFile.Name(), startTime)

	util.LogInfo("")
	util.LogInfo("Finished build (#" + strconv.Itoa(buildId) + ") in " + strconv.Itoa(b.BuildingTime) + "ms")
	util.LogInfo("")
	util.LogInfo("------------------------------------------------------")
	return nil, b
}

func createWorkingDir() error {
	util.LogInfo("Creating working directory (./work/)")
	err := os.Mkdir("work", os.ModePerm)
	return err
}

func buildProject(buildScript string) (error, bool) {
	util.LogInfo("Building the project (" + buildScript + ")")
	command := exec.Command(buildScript)
	// command.Stdout = os.Stdout

	err := command.Run()
	if err != nil {
		return err, false
	}

	success := command.ProcessState.ExitCode() == 0
	return nil, success
}

func getArtifactFile(path string, buildId int) (*os.File, error) {
	util.LogInfo("Getting the artifact file")

	matches, err := filepath.Glob(path)
	if err != nil {
		return nil, errors.New("could not find artifact file")
	}

	if len(matches) == 0 {
		return nil, errors.New("could not find artifact file")
	}

	file, err := os.OpenFile(matches[0], 0, os.ModePerm)
	if err != nil {
		return nil, errors.New("could not find artifact file")
	}

	defer file.Close()

	stat, _ := file.Stat()
	fileName := stat.Name()

	// copy file
	buildDir := "builds/"
	err = os.Mkdir(buildDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	buildDir = buildDir + "build-" + strconv.Itoa(buildId) + "/"
	err = os.Mkdir(buildDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	copiedFile, err := os.Create(buildDir + fileName)
	if err != nil {
		return nil, errors.New("could not copy artifact file")
	}

	err = util.FastCopyFile(file, copiedFile)
	if err != nil {
		return nil, errors.New("could not copy artifact file")
	}

	return copiedFile, err
}

func getGitInfo() (error, string, string) {
	util.LogInfo("Getting information about the latest commit")
	gitApp, err := git.PlainOpen("work")
	if err != nil {
		util.LogInfo(err.Error())
		return err, "", ""
	}

	head, err := gitApp.Head()
	if err != nil {
		util.LogInfo(err.Error())
		return err, "", ""
	}

	hash := head.Hash()

	commit, err := gitApp.CommitObject(hash)
	if err != nil {
		util.LogInfo(err.Error())
		return err, "", ""
	}

	message := commit.Message
	message = message[0 : len(message)-1] // remove last '\n'

	return nil, hash.String(), message
}
