package controllers

import (
	"context"
	"net/http"
	"superPeer/utils"

	"github.com/gin-gonic/gin"
)

func PublishTrackController(c *gin.Context) {
    // Define a struct that matches your incoming JSON
    type RequestBody struct {
        TrackId        string `json:"trackId"`
        PublisherName string `json:"publisherName"`
        Size          int    `json:"size"`  // Changed from string to int
    }

    var requestBody RequestBody
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    ctx := context.Background()
    redisClient := utils.RedisClient

    // Store track metadata
    _, err := redisClient.HSet(ctx, "track:"+requestBody.TrackId,
        "publisherName", requestBody.PublisherName,
        "size", requestBody.Size,  // Now passing an int instead of string
        "peerAvailable", false,
    ).Result()

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish track"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "trackId": requestBody.TrackId,
        "statusMessage": "Track published successfully",
    })
}