package services

import (
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
    collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
    return &UserService{collection: collection}
}

func (s *UserService) CreateUser(user models.User) (primitive.ObjectID, error) {
    result, err := s.collection.InsertOne(context.Background(), user)
    if err != nil {
        return primitive.NilObjectID, err
    }
    return result.InsertedID.(primitive.ObjectID), nil
}

func (s *UserService) GetUser(id primitive.ObjectID) (*models.User, error) {
    var user models.User
    err := s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *UserService) UpdateUser(id primitive.ObjectID, updateData bson.M) (*mongo.UpdateResult, error) {
    filter := bson.M{"_id": id}
    update := bson.M{"$set": updateData}
    return s.collection.UpdateOne(context.Background(), filter, update)
}

func (s *UserService) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
    filter := bson.M{"_id": id}
    return s.collection.DeleteOne(context.Background(), filter)
}

func (s *UserService) GetUserCount() (int64, error) {
    return s.collection.CountDocuments(context.Background(), bson.M{})
}

func (s *UserService) ListUser() ([]models.User, error) {
    var users []models.User
    cursor, err := s.collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    for cursor.Next(context.Background()) {
        var user models.User
        if err := cursor.Decode(&user); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}
