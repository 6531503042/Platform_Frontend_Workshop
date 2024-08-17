package main

import (
	"backend/routes"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// Initialize Redis and MongoDB
	utils.InitRedis()
	utils.InitMongoDB()

	// Initialize Fiber application
	app := fiber.New()

	// Set up routes
app.Use(cors.New(cors.Config{
	AllowOrigins: "http://localhost:3000",
	AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",

}))

	routes.Setup(app)

	// Start server on port 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
