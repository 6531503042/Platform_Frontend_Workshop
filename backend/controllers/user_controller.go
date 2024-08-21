package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	service *services.UserService
}

func NewUserController() *UserController {
	userCollection := utils.MongoDB.Collection("users")
	return &UserController{
		service: services.NewUserService(userCollection),
	}
}

// CreateUser handles user creation requests.
func (uc *UserController) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload", "details": err.Error()})
	}

	result, err := uc.service.CreateUser(user)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user", "details": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetUser handles fetching a user by ID.
func (uc *UserController) GetUser(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	user, err := uc.service.GetUser(id)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user", "details": err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UpdateUser handles requests to update a user's information.
func (uc *UserController) UpdateUser(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload", "details": err.Error()})
	}

	result, err := uc.service.UpdateUser(id, updateData)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user", "details": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// DeleteUser handles requests to delete a user by ID.
func (uc *UserController) DeleteUser(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID format"})
	}

	result, err := uc.service.DeleteUser(id)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user", "details": err.Error()})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}

// GetUserCount handles requests to get the total number of users.
func (uc *UserController) GetUserCount(c *fiber.Ctx) error {
	count, err := uc.service.GetUserCount()
	if err != nil {
		log.Printf("Failed to get user count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user count", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"userCount": count})
}

// ListUsers handles requests to list all users.
func (uc *UserController) ListUsers(c *fiber.Ctx) error {
	users, err := uc.service.ListUser()
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list users", "details": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// GetUserStatistics provides aggregated data for charts.
func (uc *UserController) GetUserStatistics(c *fiber.Ctx) error {
	collection := utils.MongoDB.Collection("users")

	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$substr", Value: bson.A{"$createdAt", 0, 7}}}}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Printf("Aggregation error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error", "details": err.Error()})
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Printf("Cursor error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(results)
}
