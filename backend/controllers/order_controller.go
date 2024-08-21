package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderController struct to hold service instance
type OrderController struct {
	service *services.OrderService
}

// NewOrderController creates a new instance of OrderController
func NewOrderController() *OrderController {
	orderCollection := utils.MongoDB.Collection("orders")
	return &OrderController{
		service: services.NewOrderService(orderCollection),
	}
}

// @Summary Create an order
// @Description Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.Order true "Order"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders [post]
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

// @Summary Get an order by ID
// @Description Fetch an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id} [get]
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

// @Summary Get all orders
// @Description Retrieve a list of all orders
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]interface{}
// @Router /orders [get]
func (oc *OrderController) GetAllOrders(c *fiber.Ctx) error {
    orders, err := oc.service.GetAllOrders()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(orders)
}

// @Summary Update an order by ID
// @Description Modify the details of an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param updateData body map[string]interface{} true "Order Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id} [put]
func (oc *OrderController) UpdateOrder(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
    }

    var update map[string]interface{}
    if err := c.BodyParser(&update); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    result, err := oc.service.UpdateOrder(id, update)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary Delete an order by ID
// @Description Remove an order by its ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/{id} [delete]
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


// @Summary Get order statistics
// @Description Retrieve aggregated statistics for orders
// @Tags Orders
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /orders/statistics [get]
func (oc *OrderController) GetOrderStatistics(c *fiber.Ctx) error {
    statistics, err := oc.service.GetOrderStatistics()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to get order statistics",
            "details": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(statistics)
}