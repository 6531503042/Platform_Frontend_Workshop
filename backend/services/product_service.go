package services

import (
	"backend/models"
	"backend/utils"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProduct (product models.Product) (*mongo.InsertOneResult, error) {
	collection := utils.MongoDB.Collection("products")

	result, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		return nil, err
	}

	return result, nil
}


// func GetProduct(id primitive.ObjectID) (*models.Product, error) {
// 	collec
// }