package main

import (
	"careerconnect/backend/database"
	"careerconnect/backend/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	r := gin.Default()
	r.POST("/register", handlers.RegisterUser)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}