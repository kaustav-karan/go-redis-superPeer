package controllers

import (
	"context"
	"net/http"
	// "superPeer/models"
	"superPeer/utils"

	"github.com/gin-gonic/gin"
)

func FetchTrackUriController(c *gin.Context) {
	var request struct {
		TrackId string `json:"trackId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx := context.Background()
	redisClient := utils.RedisClient

	// Check if track exists
	exists, err := redisClient.Exists(ctx, "track:"+request.TrackId).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Track not found"})
		return
	}

	// Get peerAvailable status
	peerAvailable, err := redisClient.HGet(ctx, "track:"+request.TrackId, "peerAvailable").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if peerAvailable == "1" {
		// Get trackUri from Redis
		trackUri, err := redisClient.HGet(ctx, "track:"+request.TrackId, "trackUri").Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		// Return the trackUri
		// Check if trackUri is empty or not
		if trackUri == "" {
			currentServerIp := utils.GetServerIP()
			// If trackUri is empty, return the current server IP
			trackUri = currentServerIp
		}
		c.JSON(http.StatusOK, gin.H{"trackUri": trackUri})
	} else {
		// Return default server IP (should be from config)
		currentServerIp := utils.GetServerIP()
		c.JSON(http.StatusOK, gin.H{"trackUri": currentServerIp})
	}
}