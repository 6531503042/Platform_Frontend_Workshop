package main

import (
	"backend/routes"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Initialize Redis and MongoDB
	utils.InitRedis()
	utils.InitMongoDB()

	// Initialize Fiber application
	app := fiber.New()

	// Set up routes
	routes.Setup(app)

	// Start server on port 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
