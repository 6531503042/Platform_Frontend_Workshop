package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(c * fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}

	result, err := services.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

func GetUser(c * fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user, err := services.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user")
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}


// UpdateUser handles requests to update a user's information.
func UpdateUser(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := services.UpdateUser(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// DeleteUser handles requests to delete a user by ID.
func DeleteUser(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	result, err := services.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.Status(fiber.StatusOK).SendString("User deleted successfully")
}

// GetUserCount handles requests to get the total number of users.
func GetUserCount(c *fiber.Ctx) error {
	count, err := services.GetUserCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"userCount": count})
}

func ListUsers(c *fiber.Ctx) error {
	users, err := services.ListUser()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// GetUserStatistics provides aggregated data for charts.
func GetUserStatistics(c *fiber.Ctx) error {
	collection := utils.MongoDB.Collection("users")

	// Example: Get user count by month
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{{"_id", bson.D{{"$substr", bson.A{"$createdAt", 0, 7}}}}, {"count", bson.D{{"$sum", 1}}}}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(results)
}