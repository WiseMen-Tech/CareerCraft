package handlers

import (
	"fmt"
	"os"
	"context"
	"net/http"
	"time"

	"careerconnect/backend/database"
	"careerconnect/backend/utils"
	"careerconnect/backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles POST /register
func RegisterUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ✅ Check if user already exists
	var existing models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error securing password"})
		return
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	_, err = database.UserCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ User registered successfully", "userId": user.ID})
}

// LoginUser handles POST /login
func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find user by email
	var user models.User
	err := database.UserCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Return token
	c.JSON(http.StatusOK, gin.H{
		"message": "✅ Login successful",
		"token":   token,
	})
}
// CreateProfile handles POST /profile
func CreateProfile(c *gin.Context) {
	// Get userId from token
	userIdVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId, _ := primitive.ObjectIDFromHex(userIdVal.(string))

	// Parse form-data (not JSON)
	education := c.PostForm("education")
	location := c.PostForm("location")
	phone := c.PostForm("phone")
	skills := c.PostFormArray("skills")
	interests := c.PostFormArray("interests")

	// Handle optional resume file
	var resumes []string
	file, err := c.FormFile("resume")
	if err == nil { // resume was uploaded
		uploadDir := "./uploads"
		_ = os.MkdirAll(uploadDir, os.ModePerm)

		resumePath := fmt.Sprintf("%s/%s_%s", uploadDir, userId.Hex(), file.Filename)
		if err := c.SaveUploadedFile(file, resumePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving resume"})
			return
		}
		resumes = append(resumes, resumePath)
	}

	// Create profile object
	profile := models.UserProfile{
		ID:        primitive.NewObjectID(),
		UserID:    userId,
		Education: education,
		Skills:    skills,
		Interests: interests,
		Location:  location,
		Phone:     phone,
		ResumeURL: resumes, 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = database.ProfileCollection.InsertOne(ctx, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Profile created successfully", "resume": profile.ResumeURL})
}

func GetMyProfile(c *gin.Context) {
	userIdVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId, _ := primitive.ObjectIDFromHex(userIdVal.(string))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var profile models.UserProfile
	err := database.ProfileCollection.FindOne(ctx, bson.M{"userId": userId}).Decode(&profile)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}
// UpdateProfile handles PUT /users/profile
func UpdateProfile(c *gin.Context) {
	userIdVal, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId, _ := primitive.ObjectIDFromHex(userIdVal.(string))

	var updateData models.UserProfile
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"education": updateData.Education,
			"skills":    updateData.Skills,
			"interests": updateData.Interests,
			"location":  updateData.Location,
			"phone":     updateData.Phone,
		},
	}

	_, err := database.ProfileCollection.UpdateOne(ctx, bson.M{"userId": userId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Profile updated successfully"})
}
// UploadResume handles POST /users/resumes
func UploadResume(c *gin.Context) {
	userIdVal, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdVal.(string))

	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	uploadDir := "./uploads"
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	// Store uniquely named file
	filePath := fmt.Sprintf("%s/%s_%s", uploadDir, userId.Hex(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving resume"})
		return
	}

	// Append resume to user's profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = database.ProfileCollection.UpdateOne(ctx,
		bson.M{"userId": userId},
		bson.M{"$push": bson.M{"resumes": filePath}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving resume reference"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Resume uploaded successfully", "resume": filePath})
}
// DeleteResume handles DELETE /users/resumes/:filename
func DeleteResume(c *gin.Context) {
	userIdVal, _ := c.Get("userId")
	userId, _ := primitive.ObjectIDFromHex(userIdVal.(string))

	filename := c.Param("filename")
	filePath := fmt.Sprintf("./uploads/%s_%s", userId.Hex(), filename)

	// Delete file locally
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting resume file"})
		return
	}

	// Remove from MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := database.ProfileCollection.UpdateOne(ctx,
		bson.M{"userId": userId},
		bson.M{"$pull": bson.M{"resumes": filePath}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating profile resumes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Resume deleted successfully"})
}

// LogoutUser handles POST /logout
func LogoutUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Store token in blacklist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := database.BlacklistCollection.InsertOne(ctx, bson.M{"token": tokenString})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error blacklisting token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ Logged out successfully"})
}
