package request

import (
	"net/http"

	"github.com/RingierIMU/rsb-go-lib/v9/base"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.New()
)

// StartWebServer starts server in accepting requests
func StartWebServer() {
	apiEventsHandler()
	landingCheck()

	router.Run()
}

// Main page, mainly for load balancers health check
func landingCheck() *gin.Engine {
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "I'm alright, thank you!")
	})

	return router
}

// Requests entry point
func apiEventsHandler() *gin.Engine {
	defer base.RecoverFromPanic("apiEventsHandler")

	v1 := router.Group("/api")
	route(v1)

	return router
}
