package main

import (
	"careerconnect/backend/database"
	"careerconnect/backend/handlers"
	"careerconnect/backend/middleware"

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
	r.POST("/login", handlers.LoginUser)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/profile", handlers.CreateProfile)
	auth.GET("/users", handlers.GetMyProfile)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}