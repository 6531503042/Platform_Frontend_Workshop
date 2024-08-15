package controllers

import (
	"backend/models"
	"backend/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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