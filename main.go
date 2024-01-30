package main

import (
	"expensesManage/database"
	"expensesManage/routes"

	"github.com/gofiber/fiber/v2"
	
)

func main() {
	database.ConnectWithDatabase()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	routes.Setup(app)

	app.Listen(":8000")
}
