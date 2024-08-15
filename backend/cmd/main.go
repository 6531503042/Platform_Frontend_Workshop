package main

import (
	"backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

func main() {

	//TODO: Initalize Redis, MongoDB, Kafka
	utils.InitRedis()
	utils.InitMongoDB()
	utils.InitKafka()

	//TODO: Fiber create application
	app := fiber.New()

	//TODO: Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//TODO: Start server
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}


}