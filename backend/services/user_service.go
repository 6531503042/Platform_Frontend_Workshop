package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	collection := utils.MongoDB.Collection("users")

	// Insert user into MongoDB
	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	// Publish to Kafka
	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to marshal user for Kafka: %v", err)
		return result, nil
	}

	msg := kafka.Message{
		Value: data,
	}
	err = utils.KafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
	}

	return result, nil
}

func GetUser (id primitive.ObjectID) (*models.User, error) {

	// Check redis cache
	cacheKey := "user:" + id.Hex()
	cachedUser, err := utils.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var user models.User
		err := json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	//Fetch from MongoDB
	collection := utils.MongoDB.Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id":id}).Decode(&user)

	if err != nil {
		return nil, err
	}

	// Cache result in Redis
	userData, _ := json.Marshal(user)
	utils.RedisClient.Set(context.Background(), cacheKey, userData, 0)

	return &user, nil
}