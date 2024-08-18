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

func CreateProduct( c* fiber.Ctx) error {
	var product models.Product


	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := services.CreateProduct(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}


func GetProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	product, err := services.GetProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product")
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	return c.Status(fiber.StatusOK).JSON(product)
}


func DeleteProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalud product ID")
	}

	result, err := services.DeleteProduct(id)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	return c.Status(fiber.StatusOK).SendString("Product deleted successfully")
}


func UpdateProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := services.UpdateProduct(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update product")
	}

	return c.Status(fiber.StatusOK).JSON(result)
}


func GetProductCount(c *fiber.Ctx) error {
	count, err := services.GetProductCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product count")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"productCount": count})
}

func ListProduct(c *fiber.Ctx) error {
	products, err := services.ListProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get products")
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func GetProductStatistics (c * fiber.Ctx) error {
	collection := utils.MongoDB.Collection("products")

	//Pipeline
	pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "$substr", Value: bson.A{"$createdAt", 0, 7}}}}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
        {{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
    }

	results, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer results.Close(context.Background())

	var statistics []bson.M

	if err = results.All(context.Background(), &statistics); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(statistics)
}