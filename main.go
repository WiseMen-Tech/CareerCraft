package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"CareerConnect/handlers"

	"github.com/robfig/cron/v3"
)

func main() {
	router := gin.Default()

	// Health Check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "SIH 2025 Backend running ✅"})
	})

	// Routes
	router.POST("/candidate", handlers.RecommendHandler)

	// Start cron job for alerts
	c := cron.New()
	c.AddFunc("@daily", func() {
		fmt.Println("🔔 Running daily deadline check...")
		handlers.CheckDeadlines()
	})
	c.Start()

	// Start server
	router.Run(":8080")
}
