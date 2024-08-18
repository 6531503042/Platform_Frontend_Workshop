package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProductController struct to hold service instance
type ProductController struct {
	service *services.ProductService
}

// NewProductController creates a new instance of ProductController
func NewProductController() *ProductController {
	// Get MongoDB and Redis clients
	productCollection := utils.MongoDB.Collection("products")
	redisClient := utils.RedisClient // Make sure you have this initialized in your utils

	return &ProductController{
		service: services.NewProductService(productCollection, redisClient),
	}
}
func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Add validation here
	if product.Name == "" || product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product data"})
	}

	result, err := pc.service.CreateProduct(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}


func (pc *ProductController) GetProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	product, err := pc.service.GetProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product")
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	result, err := pc.service.DeleteProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	return c.Status(fiber.StatusOK).SendString("Product deleted successfully")
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
	}

	var updateData bson.M
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := pc.service.UpdateProduct(id, updateData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update product")
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (pc *ProductController) GetProductCount(c *fiber.Ctx) error {
	count, err := pc.service.GetProductCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product count")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"productCount": count})
}

func (pc *ProductController) ListProduct(c *fiber.Ctx) error {
	products, err := pc.service.ListProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get products")
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (pc *ProductController) GetProductStatistics(c *fiber.Ctx) error {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$substr", Value: bson.A{"$createdAt", 0, 7}}}}, // Group by month
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}, // Count documents
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}}, // Sort by month
	}

	statistics, err := pc.service.AggregateProducts(pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(statistics)
}
