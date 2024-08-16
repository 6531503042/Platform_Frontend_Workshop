package utils

import (
	"backend/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB client and database variables
var MongoClient *mongo.Client
var MongoDB *mongo.Database
var ctx = context.Background()

// InitMongoDB initializes the MongoDB client and database
func InitMongoDB() {
	clientOptions := options.Client().ApplyURI(config.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
	MongoClient = client
	MongoDB = client.Database("test") // Replace with your database name
}
