package config

import (
	"Buildify/util"
	"encoding/json"
	"os"
	"path"
)

const (
	DefaultPort            = 1337
	DefaultBuildScriptPath = "build.bat"
	DefaultArtifactPath    = "work/build/libs/FancyNpcs*.jar"
)

type Config struct {
	Port            int
	BuildScriptPath string
	ArtifactPath    string
}

func (config Config) GetArtifactFileName() string {
	return path.Base(config.ArtifactPath)
}

func (config Config) Save() error {
	decodedJson, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.Create("config.json")
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

var CurrentConfig Config

func LoadConfig() error {
	if !util.FileExists("config.json") {
		// load default config
		CurrentConfig := Config{
			Port:            DefaultPort,
			BuildScriptPath: DefaultBuildScriptPath,
			ArtifactPath:    DefaultArtifactPath,
		}

		err := CurrentConfig.Save()
		if err != nil {
			return err
		}

		return nil
	}

	// read config
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(configFile, &CurrentConfig)
	if err != nil {
		return err
	}

	changedSomething := false

	if CurrentConfig.Port == 0 {
		CurrentConfig.Port = DefaultPort
		changedSomething = true
	}

	if len(CurrentConfig.BuildScriptPath) == 0 {
		CurrentConfig.BuildScriptPath = DefaultBuildScriptPath
		changedSomething = true
	}

	if len(CurrentConfig.ArtifactPath) == 0 {
		CurrentConfig.BuildScriptPath = DefaultArtifactPath
		changedSomething = true
	}

	if changedSomething {
		err = CurrentConfig.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
