package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"log"

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
	productCollection := utils.MongoDB.Collection("products")
	redisClient := utils.RedisClient
	return &ProductController{
		service: services.NewProductService(productCollection, redisClient),
	}
}

// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags Products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product Data"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [post]
func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		log.Printf("Error parsing product data: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product data"})
	}

	if product.Name == "" || product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Product name and price must be valid"})
	}

	result, err := pc.service.CreateProduct(product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// @Summary Get a product by ID
// @Description Fetch a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [get]
func (pc *ProductController) GetProductById(c *fiber.Ctx) error {
	idHex := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		log.Printf("Invalid product ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	product, err := pc.service.GetProduct(id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		log.Printf("Error retrieving product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve product"})
	}

	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}
// @Summary Delete a product by ID
// @Description Remove a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [delete]
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

	// if result.DeletedCount == 0 {
	// 	return c.Status(fiber.StatusNotFound).SendString("Product not found")
	// }

	// return c.Status(fiber.StatusOK).JSON(DeleteResult{DeletedCount: result.DeletedCount})

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary Update a product by ID
// @Description Modify the details of a product by its ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param updateData body bson.M true "Product Data"
// @Success 200 {object} UpdateResult
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [put]
func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
    idHex := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idHex)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid product ID")
    }

    // var updateData bson.M
    // if err := c.BodyParser(&updateData); err != nil {
    //     return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    // }

	var update map[string]interface{}
    if err := c.BodyParser(&update); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    result, err := pc.service.UpdateProduct(id, update)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to update product")
    }

    return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary Get the count of products
// @Description Retrieve the total number of products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /products/count [get]
func (pc *ProductController) GetProductCount(c *fiber.Ctx) error {
	count, err := pc.service.GetProductCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get product count")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"productCount": count})
}

// @Summary List all products
// @Description Retrieve a list of all products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} fiber.Map
// @Router /products [get]
func (pc *ProductController) ListProduct(c *fiber.Ctx) error {
	products, err := pc.service.ListProduct()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get products")
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// @Summary Get product statistics
// @Description Retrieve aggregated statistics for products
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} bson.M
// @Failure 500 {object} fiber.Map
// @Router /products/statistics [get]
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
		log.Printf("Error aggregating product statistics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve product statistics"})
	}

	return c.Status(fiber.StatusOK).JSON(statistics)
}
