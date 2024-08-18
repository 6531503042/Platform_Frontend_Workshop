package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
) 

func Setup(app *fiber.App) {

	app.Post("/users", controllers.CreateUser)
	app.Get("/users/:id", controllers.GetUser)
	app.Get("/users", controllers.ListUsers)
	app.Put("/users/:id", controllers.UpdateUser)
	app.Delete("/users/:id", controllers.DeleteUser)
	app.Get("/user-count", controllers.GetUserCount)
	app.Get("/user-statistics", controllers.GetUserStatistics)

	app.Post("/products", controllers.CreateProduct)
	app.Get("/products/:id", controllers.GetProduct)
	app.Get("/products", controllers.ListProduct)
	app.Put("/products/:id", controllers.UpdateProduct)
	app.Delete("/products/:id", controllers.DeleteProduct)
	app.Get("/product-count", controllers.GetProductCount)
	app.Get("/product-statistics", controllers.GetProductStatistics)
}