package builds

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

var Builds []Build

type Build struct {
	Id               int    // unique id of the build
	Time             int64  // timestamp when the build finished
	Hash             string // the hash of the commit
	Message          string // the message of the commit
	Downloads        int    // amount of downloads
	ArtifactFilePath string // path to artifact file
	BuildingTime     int    // time it took to build in milliseconds
}

func Create(id int, buildTime int64, hash string, message string, downloads int, artifactFilePath string, startTime time.Time) *Build {
	b := Build{
		Id:               id,
		Time:             buildTime,
		Hash:             hash,
		Message:          message,
		Downloads:        downloads,
		ArtifactFilePath: artifactFilePath,
		BuildingTime:     (int)(time.Now().UnixMilli() - startTime.UnixMilli()),
	}

	Builds = append(Builds, b)

	err := SaveBuildsFile("builds/")
	if err != nil {
		log.Println("Could not save build metadata")
		return nil
	}

	return &b
}

func Delete(id int, index int) error {
	err := os.RemoveAll("builds/build-" + strconv.Itoa(id) + "/")
	if err != nil {
		log.Println("Could not delete build " + strconv.Itoa(id))
		log.Println(err)
		return err
	}

	Builds = append(Builds[:index], Builds[index+1:]...)

	err = SaveBuildsFile("builds/")
	if err != nil {
		log.Println("Could not save build metadata")
	}

	return nil
}

func SaveBuildsFile(dir string) error {
	decodedJson, err := json.MarshalIndent(Builds, "", "  ")
	if err != nil {
		return err
	}

	err = os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	file, err := os.Create(dir + "builds.json")
	if err != nil {
		return err
	}

	_, err = file.Write(decodedJson)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func LoadBuildsFile(dir string) error {
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	file, err := os.ReadFile(dir + "builds.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &Builds)
	if err != nil {
		return err
	}

	return nil
}
