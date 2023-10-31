package handlers

import (
	"camping-backend/database"
	"camping-backend/enums"
	"camping-backend/models"
	"camping-backend/serializers"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func CreateCamping(c *fiber.Ctx) error {
	user, err := authUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	var request = models.Camping{
		CreatedAt: database.DB.NowFunc(),
		UpdatedAt: database.DB.NowFunc(),
		UserID:    user.ID,
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't parse createCamping",
			"data":    err.Error(),
		})
	}

	// enum check
	if err := setEnumView(request.View); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("%s", err.Error()),
		})
	}

	if err := setEnumStatus(request.IsEvCharge, request.IsSideParking, request.IsPetFriendly); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": fmt.Sprintf("%s", err.Error()),
		})
	}

	database.DB.Create(&request)

	responseUser := serializers.UserSerializer(user)
	responseCamping := serializers.CampingSerializer(&request, responseUser)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    responseCamping,
	})
}

func setEnumView(view enums.ViewKind) error {
	err := view.String()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", view, err.Error()))
	}
	return nil

}

func setEnumStatus(status ...enums.Status) error {
	for _, s := range status {
		err := s.String()
		if err != nil {
			return errors.New(fmt.Sprintf("%s: %s", s, err.Error()))
		}
	}
	return nil
}
func ListCamping(c *fiber.Ctx) error {
	var campings []models.Camping
	var owner models.User

	database.DB.Find(&campings)

	var responseCampings []serializers.Camping

	for _, camping := range campings {
		if err := FindUserById(&owner, int(camping.UserID)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Couldn't change password",
				"data":    err.Error(),
			})
		}
		responseUser := serializers.UserSerializer(&owner)
		responseCamping := serializers.CampingSerializer(&camping, responseUser)
		responseCampings = append(responseCampings, responseCamping)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    responseCampings,
	})

}
