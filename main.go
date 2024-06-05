package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mike-kimani/goauth/database"
	"github.com/mike-kimani/goauth/routes"
	"log"
)

func main() {
	database.Connect()
	app := fiber.New()

	routes.Setup(app)

	err := app.Listen(":8000")

	if err != nil {
		log.Fatal(err)
	}
}
