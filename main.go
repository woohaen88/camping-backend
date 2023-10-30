package main

import (
	"camping-backend/database"
	"camping-backend/handlers"
	"camping-backend/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

func setupRoute(app *fiber.App) {
	api := app.Group("/api/v1")

	user := api.Group("/user")
	user.Post("/", handlers.CreateUser)
	user.Post("/login", handlers.Login)
	user.Put("/change-password", middleware.JwtMiddleWare(), handlers.ChangePassword)
	user.Get("/me", middleware.JwtMiddleWare(), handlers.Me)
}

func main() {
	database.ConnectDB()
	app := fiber.New()

	setupRoute(app)

	log.Fatal(app.Listen(":3000"))

}
