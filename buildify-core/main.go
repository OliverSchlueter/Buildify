package main

import (
	"Buildify/api"
	"Buildify/builds"
	"Buildify/util"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
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

	admin := getAdminUser()

	// auto save
	go func() {
		for {
			log.Println("Saving builds.json")
			builds.SaveBuildsFile("builds/")
			time.Sleep(time.Minute * 5)
		}
	}()

	api.Start(*port, admin)
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
