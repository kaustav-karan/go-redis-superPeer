package controllers

import (
	"context"
	"net/http"
	"superPeer/utils"

	"github.com/gin-gonic/gin"
)

func NotifyNewPeerController(c *gin.Context) {
	var request struct {
		TrackId  string `json:"trackId" binding:"required"`
		PeerAvailable bool `json:"peerAvailable" binding:"required"`
		TrackUri string `json:"clientIp" binding:"required"`
	}

	// Parse the JSON body
	if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "error": "Invalid request",
        "details": err.Error(),  // This will show exactly what's wrong
    })
    return
	}

	ctx := context.Background()
	redisClient := utils.RedisClient
	if redisClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis client is not initialized"})
		return
	}

	// Use a Redis transaction (pipeline) for atomic operations
	pipe := redisClient.TxPipeline()

	// Check if the trackId already exists in Redis
	existsCmd := pipe.Exists(ctx, "track:"+request.TrackId)
	existsCmdResult, err := existsCmd.Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check track existence"})
		return
	}
	if existsCmdResult == 0 {
		// Track does not exist, create a new entry
		pipe.HSet(ctx, "track:"+request.TrackId,
			"trackUri", request.TrackUri,
			"peerAvailable", request.PeerAvailable,
		)
	} else {
		// Track exists, update only if there is new data
		existingData, err := redisClient.HGetAll(ctx, "track:"+request.TrackId).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing track data"})
			return
		}

		if existingData["trackUri"] != request.TrackUri {
			pipe.HSet(ctx, "track:"+request.TrackId, "trackUri", request.TrackUri)
		}
		if existingData["peerAvailable"] != utils.BoolToString(request.PeerAvailable) {
			pipe.HSet(ctx, "track:"+request.TrackId, "peerAvailable", request.PeerAvailable)
		}
	}

	// Execute the pipeline
	_, err = pipe.Exec(ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register track"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track registered successfully"})
}