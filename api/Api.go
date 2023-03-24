package api

import (
	"Buildify/builds"
	"Buildify/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

type BuildApi struct {
	Id           int
	Time         int64
	Hash         string
	Message      string
	DownloadLink string
}

func toApiStruct(build builds.Build) BuildApi {
	return BuildApi{
		Id:           build.Id,
		Time:         build.Time,
		Hash:         build.Hash,
		Message:      build.Message,
		DownloadLink: "/download?id=" + strconv.Itoa(build.Id),
	}
}

func getBuildList() []BuildApi {
	var apiBuilds []BuildApi
	for _, build := range builds.Builds {
		apiBuilds = append(apiBuilds, toApiStruct(build))
	}

	return apiBuilds
}

func Start(port int, admin util.AuthUser) {
	router := gin.Default()
	router.GET("/builds", apiBuilds)
	router.GET("/build", apiBuild)
	router.GET("/download", apiDownload)

	startDownload := router.Group("/startBuild", gin.BasicAuth(map[string]string{
		admin.Username: admin.Password,
	}))

	startDownload.GET("/", apiStartBuild)

	err := router.Run("localhost:" + strconv.Itoa(port))
	if err != nil {
		log.Fatal("Could not start REST API")
	}
}

func apiBuilds(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, getBuildList())
}

func apiBuild(context *gin.Context) {
	idStr := context.Request.URL.Query().Get("id")

	if idStr == "latest" {
		context.IndentedJSON(http.StatusOK, toApiStruct(builds.Builds[len(builds.Builds)-1]))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	for _, build := range builds.Builds {
		if build.Id == id {
			context.IndentedJSON(http.StatusOK, toApiStruct(build))
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, "Could not find build")
}

func apiStartBuild(context *gin.Context) {
	build := builds.BuildBuild(util.BuildScriptPath, util.ResultPath)

	context.IndentedJSON(http.StatusOK, toApiStruct(build))
}

func apiDownload(context *gin.Context) {
	idStr := context.Request.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	for _, build := range builds.Builds {
		if build.Id == id {
			file, err := os.Open(build.ResultFilePath)
			if err != nil {
				return
			}

			stat, err := file.Stat()
			if err != nil {
				return
			}

			context.FileAttachment(build.ResultFilePath, stat.Name())
			return
		}
	}
}
