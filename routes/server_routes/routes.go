package server_routes

import (
	"superPeer/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Super Peer is up and running 🚀")
	})
	
	router.POST("/publishTrack", controllers.PublishTrackController)
}