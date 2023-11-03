package handlers

import (
	"camping-backend/common/errors"
	commonError "camping-backend/common/errors"
	"camping-backend/database"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"github.com/gofiber/fiber/v2"
)

func DeleteAmenity(c *fiber.Ctx) error {
	amenityId, err := c.ParamsInt("amenityId")

	authUser, _ := middleware.GetAuthUser(c)

	if err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	var amenity models.Amenity
	if err := database.DB.First(&amenity, "id = ?", amenityId).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
	}

	if amenity.UserId != authUser.ID {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, nil, "남에껄 삭제하려하면 오또케~")
	}

	if err := database.DB.Delete(&amenity).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success delete tag",
		"data":    nil,
	})
}

func UpdateAmenity(c *fiber.Ctx) error {
	amenityId, err := c.ParamsInt("amenityId")

	authUser, _ := middleware.GetAuthUser(c)

	if err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	var amenity models.Amenity
	if err := database.DB.First(&amenity, "id = ?", amenityId).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
	}

	if amenity.UserId != authUser.ID {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, nil, "남에껄 수정하려하면 오또케~")
	}

	request := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}
	if err := c.BodyParser(&request); err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	amenity.UpdatedAt = database.DB.NowFunc()
	if err := database.DB.Model(&amenity).Updates(request).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	serializedAmenity := serializers.AmenitySerializer(amenity, serializers.UserSerializer(authUser))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "update success",
		"data":    serializedAmenity,
	})
}

func ListAmenity(c *fiber.Ctx) error {

	var amenities []models.Amenity
	database.DB.Find(&amenities)

	var serializedAmenities []serializers.Amenity
	for _, amenity := range amenities {
		var amenityCreator models.User
		if err := FindUserById(&amenityCreator, int(amenity.UserId)); err != nil {
			return errors.ErrorHandler(c, fiber.StatusNotFound, err)
		}
		serializedAmenity := serializers.AmenitySerializer(amenity, serializers.UserSerializer(&amenityCreator))
		serializedAmenities = append(serializedAmenities, serializedAmenity)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedAmenities,
	})
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

	if err := database.DB.Create(&amenity).Error; err != nil {
		return errors.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	serializedUser := serializers.UserSerializer(authUser)
	serializedAmenity := serializers.AmenitySerializer(amenity, serializedUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedAmenity,
	})

}
