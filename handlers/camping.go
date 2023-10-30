package handlers

import (
	"camping-backend/database"
	"camping-backend/models"
	"camping-backend/serializers"
	"github.com/gofiber/fiber/v2"
	"time"
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

	database.DB.Create(&request)

	type ResponseUser = struct {
		ID       uint
		Email    string
		Username string
	}

	responseCamping := struct {
		ID        uint
		Title     string
		Address   string
		CreatedAt time.Time
		UpdatedAt time.Time
		User      ResponseUser
	}{
		ID:        request.ID,
		Title:     request.Title,
		Address:   request.Address,
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
		User: struct {
			ID       uint
			Email    string
			Username string
		}{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    responseCamping,
	})
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
		responseUser := serializers.UserSerializer(owner)
		responseCamping := serializers.CampingSerializer(camping, responseUser)
		responseCampings = append(responseCampings, responseCamping)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "success",
		"data":    responseCampings,
	})

}
