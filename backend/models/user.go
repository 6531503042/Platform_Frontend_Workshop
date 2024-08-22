package models

type User struct {
	// ID primitive.ObjectID `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}