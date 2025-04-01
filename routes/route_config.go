package routes

import (
	"github.com/gin-gonic/gin"
	"superPeer/routes/peer_routes"
	"superPeer/routes/root_routes"
	"superPeer/routes/server_routes"
)

func SetupRoutes(router *gin.Engine) {
	// Root routes
	root_routes.RegisterRoutes(router)
	
	// Peer routes with "/peer" prefix
	peerGroup := router.Group("/peer")
	peer_routes.RegisterRoutes(peerGroup)
	
	// Server routes with "/server" prefix
	serverGroup := router.Group("/server")
	server_routes.RegisterRoutes(serverGroup)
}