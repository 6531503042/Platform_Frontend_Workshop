package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity int `json:"quantity" bson:"quantity"`
	Status bool `json:"status" bson:"status"`
}