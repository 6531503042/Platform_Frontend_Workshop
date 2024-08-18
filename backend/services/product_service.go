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

func CreateProduct (product models.Product) (*mongo.InsertOneResult, error) {
	collection := utils.MongoDB.Collection("products")

	result, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}

	// Cache result in redis
	productData, _ := json.Marshal(product)
	cacheKey := "product:" + product.ID.Hex()
	utils.RedisClient.Set(context.Background(), cacheKey, productData, 0)

	return result, nil
}

func GetProduct(id primitive.ObjectID) (*models.Product, error) {
	// Check redis cache
	cacheKey := "product:" + id.Hex()
	cachedProduct, err := utils.RedisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var product models.Product
		err := json.Unmarshal([]byte(cachedProduct), &product)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

	// Fetch from MongoDB
	collection := utils.MongoDB.Collection("products")
	var product models.Product
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)

	if err != nil {
		return nil, err
	}

	// Cache result in Redis
	productData, _ := json.Marshal(product)
	utils.RedisClient.Set(context.Background(), cacheKey, productData, 0)

	return &product, nil
}


func ListProduct() ([]models.Product, error) {
	collection := utils.MongoDB.Collection("products")
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	var products []models.Product

	if err = cursor.All(context.Background(), &products); err != nil {
		return nil, err
	}
	return products, nil
}


func UpdateProduct(id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error) {
	collection := utils.MongoDB.Collection("products")
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Invalidate the cache
	cacheKey := "product:" + id.Hex()
	utils.RedisClient.Del(context.Background(), cacheKey)

	return result, nil
}


func DeleteProduct(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection := utils.MongoDB.Collection("products")
	cacheKey := "product:" + id.Hex()

	// Delete product from MongoDB
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	// Invalidate the cache
	utils.RedisClient.Del(context.Background(), cacheKey)

	return result, nil
}

func GetProductCount() (int64, error) {
	collection := utils.MongoDB.Collection("products")

	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}