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

}