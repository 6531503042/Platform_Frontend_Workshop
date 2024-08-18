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

// CreateUser handles user creation requests.
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := services.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetUser handles fetching a user by ID.
func GetUser(c *fiber.Ctx) error {
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

// ListUsers handles requests to list all users.
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

    // Define the aggregation pipeline
    pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$substr", Value: bson.A{"$createdAt", 0, 7}}}}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
        {{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
    }

    // Execute the aggregation
    cursor, err := collection.Aggregate(context.Background(), pipeline)
    if err != nil {
        log.Printf("Aggregation error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
    }
    defer cursor.Close(context.Background())

    // Retrieve results
    var results []bson.M
    if err := cursor.All(context.Background(), &results); err != nil {
        log.Printf("Cursor error: %v", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
    }

    // Return results
    return c.Status(fiber.StatusOK).JSON(results)
}
