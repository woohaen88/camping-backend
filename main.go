package main

import (
	"camping-backend/database"
	"camping-backend/enums"
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
	user.Put("/change-password", middleware.Protected(), handlers.ChangePassword)
	user.Get("/me", middleware.Protected(), handlers.Me)

	camping := api.Group("/camping")
	camping.Post("/", middleware.Protected(), handlers.CreateCamping)
	camping.Get("/", handlers.ListCamping)
	camping.Get("/:campingId", handlers.GetCamping)
	camping.Put("/:campingId", middleware.Protected(), handlers.UpdateCamping)
	camping.Delete("/:campingId", middleware.Protected(), handlers.DeleteCamping)

	tag := api.Group("/tag", middleware.Protected(), middleware.AssignRole(enums.Admin))
	tag.Get("/", handlers.ListTag)
	tag.Post("/", handlers.CreateTag)
	tag.Put("/:tagId", handlers.UpdateTag)
	tag.Delete("/:tagId", handlers.DeleteTag)

	amenity := api.Group("/amenity", middleware.Protected(), middleware.AssignRole(enums.Admin))
	amenity.Get("/", handlers.ListAmenity)
	amenity.Post("/", handlers.CreateAmenity)
	amenity.Put("/:amenityId", handlers.UpdateAmenity)
	amenity.Delete("/:amenityId", handlers.DeleteAmenity)

}

func main() {

	database.ConnectDB()
	app := fiber.New()

	setupRoute(app)

	log.Fatal(app.Listen(":3000"))

}
