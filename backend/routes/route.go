package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
) 

func Setup(app *fiber.App) {

	app.Post("/users", controllers.CreateUser)
	app.Get("/users/:id", controllers.GetUser)
}