package handlers

import (
	"camping-backend/common/errors"
	"camping-backend/database"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"github.com/gofiber/fiber/v2"
)

func ListAmenity(c *fiber.Ctx) error {
	return nil
}

func CreateAmenity(c *fiber.Ctx) error {
	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return errors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	request := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}

	if err := c.BodyParser(&request); err != nil {
		return errors.ErrorHandler(c, fiber.StatusBadRequest, err, "could not parse body")
	}

	var amenity models.Amenity
	amenity.Name = request.Name
	amenity.Description = request.Description
	amenity.CreatedAt = database.DB.NowFunc()
	amenity.UpdatedAt = database.DB.NowFunc()
	amenity.UserId = authUser.ID

	database.DB.Create(&amenity)

	// TODO: serializer

	serializedUser := serializers.UserSerializer(authUser)
	serializedAmenity := serializers.AmenitySerializer(amenity, serializedUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedAmenity,
	})

}
