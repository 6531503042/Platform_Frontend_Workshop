package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"encoding/json"

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

	// Cache result in redis
	userData, _ := json.Marshal(user)
	cacheKey := "user:" + user.ID.Hex()
	utils.RedisClient.Set(context.Background(), cacheKey, userData, 0)

	return result, nil
}

func GetUser(id primitive.ObjectID) (*models.User, error) {

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

	// Fetch from MongoDB
	collection := utils.MongoDB.Collection("users")
	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)

	if err != nil {
		return nil, err
	}

	// Cache result in Redis
	userData, _ := json.Marshal(user)
	utils.RedisClient.Set(context.Background(), cacheKey, userData, 0)

	return &user, nil
}


func ListUser() ([]models.User, error) {
	collection := utils.MongoDB.Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var users []models.User

	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}



func UpdateUser(id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error) {
	collection := utils.MongoDB.Collection("users")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Invalidate the cache
	cacheKey := "user:" + id.Hex()
	utils.RedisClient.Del(context.Background(), cacheKey)

	return result, nil
}
// DeleteUser deletes a user by ID from MongoDB and Redis.
func DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := utils.MongoDB.Collection("users")
	cacheKey := "user:" + id.Hex()

	// Delete user from MongoDB
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	// Delete user from Redis cache
	utils.RedisClient.Del(context.Background(), cacheKey)

	return result, nil
}

// GetUserCount retrieves the total number of users.
func GetUserCount() (int64, error) {
	collection := utils.MongoDB.Collection("users")
	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}