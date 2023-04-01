package builds

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

var Builds []Build

type Build struct {
	Id             int
	Time           int64
	Hash           string
	Message        string
	ResultFilePath string
}

func Create(id int, time int64, hash string, message string, resultFilePath string) *Build {
	b := Build{
		Id:             id,
		Time:           time,
		Hash:           hash,
		Message:        message,
		ResultFilePath: resultFilePath,
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
