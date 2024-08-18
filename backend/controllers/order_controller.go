package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderController struct {
	service *services.OrderService
}

func NewOrderController() *OrderController {
	// Pass the MongoDB collection to the NewOrderService function
	orderCollection := utils.MongoDB.Collection("orders")
	return &OrderController{
		service: services.NewOrderService(orderCollection),
	}
}

func (oc *OrderController) CreateOrder(c *fiber.Ctx) error {
	order := new(models.Order)

	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := oc.service.CreateOrder(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (oc *OrderController) GetOrderById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	order, err := oc.service.GetOrderById(id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func (oc *OrderController) GetAllOrders(c *fiber.Ctx) error {
	orders, err := oc.service.GetAllOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func (oc *OrderController) UpdateOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var update bson.M
	if err := c.BodyParser(&update); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := oc.service.UpdateOrder(id, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (oc *OrderController) DeleteOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	result, err := oc.service.DeleteOrder(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (oc *OrderController) GetOrderStatistics(c *fiber.Ctx) error {
	// Call the GetOrderStatistics method from the service
	statistics, err := oc.service.GetOrderStatistics()
	if err != nil {
		// Return a 500 status code if there's an error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get order statistics",
			"details": err.Error(),
		})
	}

	// Return the statistics with a 200 status code
	return c.Status(fiber.StatusOK).JSON(statistics)
}