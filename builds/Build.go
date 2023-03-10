package builds

import (
	"encoding/json"
	"log"
	"os"
)

var Builds []Build

type Build struct {
	Id      int
	Hash    string
	Message string
}

func Create(id int, hash string, message string) Build {
	b := Build{
		Id:      id,
		Hash:    hash,
		Message: message,
	}

	Builds = append(Builds, b)
	return b
}

func SaveBuildsFile(dir string) error {
	decodedJson, err := json.MarshalIndent(Builds, "", "  ")
	if err != nil {
		log.Fatal(err)
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
