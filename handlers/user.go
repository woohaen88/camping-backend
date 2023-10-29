package handlers

import (
	"camping-backend/database"
	"camping-backend/models"
	"camping-backend/serializers"

	"github.com/gofiber/fiber/v2"
)



func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if userExist := checkEmailDuplicate(user); userExist {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"detail": "저기여 유저가 있어요",
		})
	}

	// password 해쉬

	database.DB.Create(user)
	responseUser := serializers.UserSerializer(*user)
	return c.Status(200).JSON(responseUser)

}

func checkEmailDuplicate(user *models.User) bool{
	database.DB.Find(user, "email = ?", user.Email)	
	return user.ID>0
}