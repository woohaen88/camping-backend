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

func DeleteTag(c *fiber.Ctx) error {
	tagId, err := c.ParamsInt("tagId")

	authUser, _ := middleware.GetAuthUser(c)

	if err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	var tag models.Tag
	if err := database.Database.Conn.First(&tag, "id = ?", tagId).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
	}

	if tag.UserId != authUser.ID {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, nil, "남에껄 삭제하려하면 오또케~")
	}

	if err := database.Database.Conn.Delete(&tag).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success delete tag",
		"data":    nil,
	})
}

func UpdateTag(c *fiber.Ctx) error {
	tagId, err := c.ParamsInt("tagId")

	authUser, _ := middleware.GetAuthUser(c)

	if err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	var tag models.Tag
	if err := database.Database.Conn.First(&tag, "id = ?", tagId).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, err)
	}

	if tag.UserId != authUser.ID {
		return commonError.ErrorHandler(c, fiber.StatusNotFound, nil, "남에껄 수정하려하면 오또케~")
	}

	request := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{}
	if err := c.BodyParser(&request); err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	tag.UpdatedAt = database.Database.Conn.NowFunc()
	if err := database.Database.Conn.Model(&tag).Updates(request).Error; err != nil {
		return commonError.ErrorHandler(c, fiber.StatusBadRequest, err)
	}

	serializedTag := serializers.TagSerializer(tag, serializers.UserSerializer(authUser))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "update success",
		"data":    serializedTag,
	})
}

func ListTag(c *fiber.Ctx) error {
	var tags []models.Tag
	database.Database.Conn.Find(&tags)

	// serializer
	var serializedTags []serializers.Tag
	for _, tag := range tags {
		var user models.User
		if err := database.Database.Conn.First(&user, tag.UserId).Error; err != nil {
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
		CreatedAt: database.Database.Conn.NowFunc(),
		UpdatedAt: database.Database.Conn.NowFunc(),
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "parse error",
			"data":    err.Error(),
		})
	}

	if err := database.Database.Conn.Create(&request).Error; err != nil {
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

func FindByTagId[T int | uint | int32](tagId T) (*models.Tag, error) {
	var tag *models.Tag

	if err := database.Database.Conn.First(&tag, "id = ?", tagId).Error; err != nil {
		return nil, err
	}

	return tag, nil
}
