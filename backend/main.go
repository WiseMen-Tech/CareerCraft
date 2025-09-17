package main

import (
	"careerconnect/backend/database"
	"careerconnect/backend/handlers"
	"careerconnect/backend/middleware"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	auth.GET("/me", handlers.GetMyProfile)
	auth.PUT("/profile", handlers.UpdateProfile)
	auth.POST("/resumes", handlers.UploadResume)
	auth.DELETE("/resumes/:filename", handlers.DeleteResume)
	auth.POST("/logout", handlers.LogoutUser)

	port := os.Getenv("PORT")
	r.Run(":" + port)
}