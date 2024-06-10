package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mike-kimani/goauth/controllers"
)

func Setup(app *fiber.App) {
	app.Get("/", controllers.Hello)
}

func Register(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
}

func Login(app *fiber.App) {
	app.Post("/api/login", controllers.Login)
}

func User(app *fiber.App) {
	app.Get("/api/user", controllers.User)
}
