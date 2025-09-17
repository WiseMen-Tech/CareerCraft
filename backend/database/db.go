package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var ProfileCollection *mongo.Collection
var BlacklistCollection *mongo.Collection

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("❌ MONGO_URI not set in .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	// check connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Could not ping MongoDB:", err)
	}

	Client = client
	UserCollection = client.Database(os.Getenv("MONGO_DB")).Collection("users")
	ProfileCollection = client.Database(os.Getenv("MONGO_DB")).Collection("profiles")
	BlacklistCollection = client.Database(os.Getenv("MONGO_DB")).Collection("blacklist")


	fmt.Println("✅ Connected to MongoDB Atlas")
}
