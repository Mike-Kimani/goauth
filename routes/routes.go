package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mike-kimani/goauth/controllers"
)

func Setup(app *fiber.App) {
	app.Get("/", controllers.Hello)
}
