package main

import (
	"log"
	"net/http"
	"os"
	
	"superPeer/routes"
	"superPeer/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start background services
	go utils.StartIPUpdater()
	go utils.StartBeaconSender()

	// Initialize Redis connection
	utils.RedisConnect()

	// Create Gin router
	router := gin.Default()
	// router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure routes
	routes.SetupRoutes(router)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server is running on http://localhost:%s 🚀", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}