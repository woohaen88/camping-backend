package handlers

import (
	"camping-backend/database"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ListTag(c *fiber.Ctx) error {

	return nil
}

func CreateTag(c *fiber.Ctx) error {

	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "couldn't load auth user",
			"data":    err.Error(),
		})
	}
	fmt.Println("authUser", authUser)

	request := models.Tag{
		UserId:    authUser.ID,
		CreatedAt: database.DB.NowFunc(),
		UpdatedAt: database.DB.NowFunc(),
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "parse error",
			"data":    err.Error(),
		})
	}

	if err := database.DB.Create(&request).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "duplicate name",
			"data":    err.Error(),
		})
	}

	serializedUser := serializers.UserSerializer(authUser)
	serializedTag := serializers.TagSerializer(request, serializedUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    serializedTag,
	})
}
