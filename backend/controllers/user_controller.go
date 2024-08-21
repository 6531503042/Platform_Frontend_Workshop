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

// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users [post]
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

// @Summary Get a user by ID
// @Description Fetch a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/{id} [get]
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

// @Summary Update a user
// @Description Update user details by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body bson.M true "User Data"
// @Success 200 {object} mongo.UpdateResult
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/{id} [put]
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


//@Summary Delete a user
//@Description Delete a user by ID
//@Tags Users
//@Accept json
//@Produce json
//@Param id path string true "User ID"
//@Success 200 {object} models.User
//@Failure 400 {object} fiber.Map
//@Failure 404 {object} fiber.Map
//@Failure 500 {object} fiber.Map
//@Router /users/{id} [delete]
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

// @Summary Get user count
// @Description Retrieve the total number of users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /users/count [get]
func (uc *UserController) GetUserCount(c *fiber.Ctx) error {
	count, err := uc.service.GetUserCount()
	if err != nil {
		log.Printf("Failed to get user count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user count", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"userCount": count})
}

// @Summary List users
// @Description Retrieve a list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} fiber.Map
// @Router /users [get]
func (uc *UserController) ListUsers(c *fiber.Ctx) error {
	users, err := uc.service.ListUser()
	if err != nil {
		log.Printf("Failed to list users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list users", "details": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// @Summary Get user statistics
// @Description Retrieve aggregated user data for charts
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} bson.M
// @Failure 500 {object} fiber.Map
// @Router /users/statistics [get]
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
