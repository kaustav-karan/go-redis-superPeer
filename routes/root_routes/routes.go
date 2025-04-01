package root_routes

import (
	"superPeer/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Super Peer is up and running 🚀")
	})
	
	router.POST("/notifyNewPeer", controllers.NotifyNewPeerController)
}