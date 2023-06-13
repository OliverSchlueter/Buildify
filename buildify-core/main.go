package main

import (
	"Buildify/api"
	"Buildify/builds"
	"Buildify/config"
	"Buildify/util"
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
)

func main() {
	util.SetStartupTime(time.Now().UnixMilli())

	// loading config
	err := config.LoadConfig()
	if err != nil {
		log.Println("Could not load config.json")
		log.Fatal(err)
		return
	}

	log.Println("Loaded config.json")
	log.Println("- Port: " + strconv.Itoa(config.CurrentConfig.Port))
	log.Println("- Build script path: " + config.CurrentConfig.BuildScriptPath)
	log.Println("- Artifact path: " + config.CurrentConfig.ArtifactPath)
	log.Println("")

	// loading admin user
	admin := getAdminUser()
	log.Println("Loaded admin user")

	// loading builds
	err = builds.LoadBuildsFile("builds/")
	if err != nil {
		log.Println("Could not load build metadata file")
	}

	log.Println("Loaded build metadata")

	// starting auto save
	go func() {
		time.Sleep(time.Minute * 2)

		for {
			log.Println("Auto saving build metadata")
			builds.SaveBuildsFile("builds/")
			time.Sleep(time.Minute * 5)
		}
	}()

	// print memory usage
	go func() {
		time.Sleep(time.Second * 5)
		for {
			util.PrintMemUsage()
			time.Sleep(time.Second * 30)
		}
	}()

	log.Println("Starting http server...")
	go api.Start(config.CurrentConfig.Port, admin)
	log.Println("Started http server")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if input == "" {
			continue
		}

		log.Println("Input: '" + input + "'")
	}
}

func getAdminUser() util.AuthUser {
	var admin util.AuthUser

	adminAuthFile, err := os.OpenFile("admin_auth.json", os.O_RDWR, os.ModePerm)
	defer adminAuthFile.Close()

	if os.IsNotExist(err) {
		uuid, err := guid.NewV4()

		admin = util.AuthUser{
			Username: "admin",
			Password: strings.Replace(uuid.String(), "-", "", -1)[:12],
		}

		adminAuthFile, err = os.Create("admin_auth.json")
		defaultAdmin, _ := json.MarshalIndent(admin, "", "\t")
		_, err = adminAuthFile.Write(defaultAdmin)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Generated default admin user (admin_auth.json)")
	} else {
		adminAuthFileContent, err := os.ReadFile("admin_auth.json")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(adminAuthFileContent, &admin)
		if err != nil {
			log.Fatal("Could not load the admin user")
		}
	}

	return admin
}
