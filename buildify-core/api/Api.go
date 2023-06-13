package api

import (
	"Buildify/builds"
	"Buildify/util"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BuildApi struct {
	Id           int
	Time         int64
	Hash         string
	Message      string
	Downloads    int
	FileName     string
	DownloadLink string
}

func toApiStruct(build builds.Build) BuildApi {
	return BuildApi{
		Id:           build.Id,
		Time:         build.Time,
		Hash:         build.Hash,
		Message:      build.Message,
		Downloads:    build.Downloads,
		FileName:     util.GetArtifactFileName(build.ArtifactFilePath),
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
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/server-stats", apiServerStats)
	router.GET("/builds", apiBuilds)
	router.GET("/build", apiBuild)
	router.GET("/download", apiDownload)

	startDownload := router.Group("/startBuild", gin.BasicAuth(map[string]string{
		admin.Username: admin.Password,
	}))

	startDownload.GET("/", apiStartBuild)

	deleteBuild := router.Group("/deleteBuild", gin.BasicAuth(map[string]string{
		admin.Username: admin.Password,
	}))

	deleteBuild.GET("/", apiDeleteBuild)

	err := router.Run("localhost:" + strconv.Itoa(port))
	if err != nil {
		log.Fatal("Could not start REST API")
	}
}

func apiServerStats(context *gin.Context) {
	mem := util.GetMemoryStats()

	context.IndentedJSON(http.StatusOK, map[string]string{
		"uptime": strconv.FormatInt(util.GetUptime(), 10),
		"memory": strconv.FormatUint(mem.Alloc, 10),
		"num-gc": strconv.FormatUint(uint64(mem.NumGC), 10),
	})
}

func apiBuilds(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin")) //TODO: I think this is not safe
	context.Header("Access-Control-Allow-Methods", "GET")
	context.IndentedJSON(http.StatusOK, getBuildList())
}

func apiBuild(context *gin.Context) {
	idStr := context.Request.URL.Query().Get("id")
	if idStr == "" {
		context.IndentedJSON(http.StatusNotFound, "Please provide a build id")
		return
	}

	if idStr == "latest" {
		context.IndentedJSON(http.StatusOK, toApiStruct(builds.Builds[len(builds.Builds)-1]))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, "Please provide a valid build id")
		return
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
	if builds.IsBuilding {
		context.IndentedJSON(http.StatusLocked, "Please wait, there is already an ongoing build")
		return
	}

	builds.IsBuilding = true
	err, build := builds.CreateBuild()
	builds.IsBuilding = false

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, "Something went wrong while building")
		return
	}

	context.IndentedJSON(http.StatusOK, toApiStruct(*build))
}

func apiDownload(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", context.GetHeader("Origin")) //TODO: I think this is not safe
	context.Header("Access-Control-Allow-Methods", "GET")

	idStr := context.Request.URL.Query().Get("id")
	if idStr == "" {
		context.IndentedJSON(http.StatusNotFound, "Please provide a build id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, "Please provide a valid build id")
		return
	}

	for i, build := range builds.Builds {
		if build.Id == id {
			file, err := os.Open(build.ArtifactFilePath)
			if err != nil {
				return
			}

			stat, err := file.Stat()
			if err != nil {
				return
			}

			context.FileAttachment(build.ArtifactFilePath, stat.Name())
			builds.Builds[i].Downloads += 1
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, "Could not find a build with this id")
}

func apiDeleteBuild(context *gin.Context) {
	idStr := context.Request.URL.Query().Get("id")
	if idStr == "" {
		context.IndentedJSON(http.StatusNotFound, "Please provide a build id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, "Please provide a valid build id")
		return
	}

	for index, build := range builds.Builds {
		if build.Id == id {

			err := builds.Delete(id, index)
			if err != nil {
				context.IndentedJSON(http.StatusInternalServerError, "Could not delete build")
				return
			}

			context.IndentedJSON(http.StatusInternalServerError, "Deleted build #"+idStr)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, "Could not find a build with this id")
}
