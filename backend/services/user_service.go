package services

import (
	"backend/models"
	"context"
	"errors"
	"fmt"

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

// FindUserByEmail finds a user by their email address
func (s *UserService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *UserService) CreateUser(user models.User) (primitive.ObjectID, error) {

    if user.Email == "" || user.Name == "" {
        return primitive.NilObjectID, errors.New("missing required fields: email and name")
    }

    // Check if the email already exists
	existingUser, err := s.FindUserByEmail(user.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return primitive.NilObjectID, fmt.Errorf("error checking email existence: %v", err)
	}
	if existingUser != nil {
		return primitive.NilObjectID, errors.New("email already exists")
	}

    // Insert the new user into the database
	result, err := s.collection.InsertOne(context.Background(), user)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to create user: %v", err)
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

    if len(updateData) == 0 {
        return nil, errors.New("no data to update")
    }

    filter := bson.M{"_id": id}
	update := bson.M{"$set": updateData}
	result, err := s.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserService) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
    filter := bson.M{"_id": id}
    result, err :=  s.collection.DeleteOne(context.Background(), filter)
    
    if err != nil {
        return nil, err
    }

    if result.DeletedCount == 0 {
        return nil, errors.New("user not found")
    }

    return result, nil
}

func (s *UserService) GetUserCount() (int64, error) {
    count, err :=  s.collection.CountDocuments(context.Background(), bson.M{})
    if err != nil {
        return 0, err
    }

    return count, nil
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
