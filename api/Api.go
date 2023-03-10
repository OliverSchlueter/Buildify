package api

import (
	"Buildify/builds"
	"Buildify/util"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Start(port int) {
	router := gin.Default()
	router.GET("/builds", apiBuilds)
	router.GET("/build", apiBuild)
	router.GET("/startBuild", apiStartBuild)

	err := router.Run("localhost:" + strconv.Itoa(port))
	if err != nil {
		log.Fatal("Could not start REST API")
	}
}

func apiBuilds(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, builds.Builds)
}

func apiBuild(context *gin.Context) {
	idStr := context.Request.URL.Query().Get("id")

	if idStr == "latest" {
		context.IndentedJSON(http.StatusOK, builds.Builds[len(builds.Builds)-1])
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}

	for _, build := range builds.Builds {
		if build.Id == id {
			context.IndentedJSON(http.StatusOK, build)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, "Could not find build")
}

func apiStartBuild(context *gin.Context) {
	build := builds.BuildBuild(util.BuildScriptPath, util.ResultPath)

	context.IndentedJSON(http.StatusOK, build)
}
