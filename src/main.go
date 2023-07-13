package main

import (
	"Buildify/api"
	"Buildify/builds"
	"Buildify/config"
	"Buildify/util"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Microsoft/go-winio/pkg/guid"
	"github.com/gin-gonic/gin"
)

func main() {
	util.SetStartupTime(time.Now().UnixMilli())

	// setting up logger
	err := os.Mkdir("logs/", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile("logs/"+time.Now().Format(time.DateOnly)+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logWriter := io.MultiWriter(os.Stdout, logFile)
	gin.DefaultWriter = logWriter
	log.SetOutput(logWriter)

	log.Println("-----------------------")
	log.Println("Starting Building...")
	log.Println("-----------------------")

	// loading config
	err = config.LoadConfig()
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

	// start http server
	go api.Start(config.CurrentConfig.Port, admin)

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	log.Println("Started http server (http://" + hostname + ":" + strconv.Itoa(config.CurrentConfig.Port) + ")")

	// start waiting for cli commands
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if input == "" {
			continue
		}

		input = strings.ToLower(input)
		parts := strings.Split(input, " ")
		cmd := parts[0]
		// args := parts[1:]

		switch cmd {
		case "help":
			log.Println("--------------------------------------")
			log.Println("help - shows this menu")
			log.Println("status - shows the server status")
			log.Println("stop - stops the application")
			log.Println("clear|cls - clears the terminal")
			log.Println("--------------------------------------")
		case "stop":
			log.Println("Stopping Buildify")
			builds.SaveBuildsFile("builds/")
			return
		case "status":
			log.Println("--------------------------------------")
			log.Println("Status: " + util.ColorGreen + "Running" + util.ColorReset)
			util.PrintUptime()
			util.PrintMemUsage()
			util.PrintAmountRequests()
			log.Println("--------------------------------------")
		case "clear":
		case "cls":
			fmt.Println("\033[2J")
		}
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
