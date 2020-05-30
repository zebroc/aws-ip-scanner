package request

import (
	"github.com/gin-gonic/gin"
)

func route(router *gin.RouterGroup) {
	router.POST("/events", func(context *gin.Context) {
		eventAction(context)
	})
}
