package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
) 

func Setup(app *fiber.App) {

	

	// Initialize controller
	userController := controllers.NewUserController()
	productController := controllers.NewProductController()
	orderController := controllers.NewOrderController()
	


	//User
	app.Post("/users", userController.CreateUser)
	app.Get("/users/:id", userController.GetUser)
	app.Get("/users", userController.ListUsers)
	app.Put("/users/:id", userController.UpdateUser)
	app.Delete("/users/:id", userController.DeleteUser)
	app.Get("/user-count", userController.GetUserCount)
	app.Get("/user-statistics", userController.GetUserStatistics)

	//Product
	app.Post("/products", productController.CreateProduct)
	app.Get("/products/:id", productController.GetProduct)
	app.Get("/products", productController.ListProduct)
	app.Put("/products/:id", productController.UpdateProduct)
	app.Delete("/products/:id", productController.DeleteProduct)
	app.Get("/product-count", productController.GetProductCount)
	app.Get("/product-statistics", productController.GetProductStatistics)

	//Order
	app.Post("/orders", orderController.CreateOrder)
	app.Get("/orders/:id", orderController.GetOrderById)
	app.Get("/orders", orderController.GetAllOrders)
	app.Put("/orders/:id", orderController.UpdateOrder)
	app.Delete("/orders/:id", orderController.DeleteOrder)
	app.Get("/order-statistics", orderController.GetOrderStatistics)
}