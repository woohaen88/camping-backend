package main

import (
	"camping-backend/database"
	"camping-backend/handlers"

	"github.com/gofiber/fiber/v2"
)



func setupRoute(app *fiber.App){
	api :=app.Group("/api/v1")

	user := api.Group("/user")
	user.Post("/", handlers.CreateUser)
}

func main(){
	database.ConnectDB()
	app := fiber.New()

	setupRoute(app)

	app.Listen(":3000")


}