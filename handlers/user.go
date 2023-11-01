package handlers

import (
	"camping-backend/database"
	"camping-backend/middleware"
	"camping-backend/models"
	"camping-backend/serializers"
	"errors"
	"fmt"
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

	if err := user.Role.Check(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "success",
			"message": "user role 타입이 올바르지 않아욧",
			"data":    err.Error(),
		})
	}

	// password 해쉬
	user.PaswordHash(user.Password)

	database.DB.Create(user)
	responseUser := serializers.UserSerializer(user)
	return c.Status(200).JSON(responseUser)

}

func Me(c *fiber.Ctx) error {
	user, err := middleware.GetAuthUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	fmt.Println("user: ", user)

	return nil
}

func ChangePassword(c *fiber.Ctx) error {
	user, err := middleware.GetAuthUser(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"detail": err.Error(),
		})
	}

	type Req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		CheckPassword   string `json:"check_password"`
	}

	req := Req{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Couldn't change password",
			"data":    err.Error(),
		})
	}

	if !CheckPasswordHash(req.CurrentPassword, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "It is different from the previous password.",
			"data":    nil,
		})
	}

	if req.NewPassword != req.CheckPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "The password you entered is different",
			"data":    nil,
		})
	}

	user.PaswordHash(req.NewPassword)

	database.DB.Save(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "The Password Change Successful",
		"data":    nil,
	})
}

func FindUserById(user *models.User, id int) error {
	database.DB.First(&user, id)
	if user.ID == 0 {
		return errors.New("해당 유저가 없어용")
	}
	return nil

}
