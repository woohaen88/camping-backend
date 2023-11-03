package handlers

import (
	commonError "camping-backend/common/errors"
	"camping-backend/database"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ListTag(c *fiber.Ctx) error {
	var tags []models.Tag
	database.DB.Find(&tags)

	// serializer
	var serializedTags []serializers.Tag
	for _, tag := range tags {
		var user models.User
		if err := database.DB.First(&user, tag.UserId).Error; err != nil {
			return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
		}
		serializedTag := serializers.TagSerializer(tag, serializers.UserSerializer(&user))
		serializedTags = append(serializedTags, serializedTag)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "error",
		"message": "error",
		"data":    serializedTags,
	})
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
