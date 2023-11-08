package handlers

import (
	"camping-backend/config"
	"camping-backend/database"
	"camping-backend/models"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := new(LoginInput)
	var userData UserData

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on login request",
			"data":    err.Error(),
		})
	}

	email := input.Email
	password := input.Password
	userModel, err := new(models.User), *new(error)

	if isEmail(email) {
		userModel, err = getUserByEmail(email)
	}

	if userModel == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	userData = UserData{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password,
	}

	if !CheckPasswordHash(password, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid password",
			"data":    nil,
		})

	}
	// jwt
	token := jwt.New(jwt.SigningMethodHS256)
	claim := token.Claims.(jwt.MapClaims)
	claim["userId"] = userData.ID
	claim["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET"))) // Todo

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data":    t,
	})
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func getUserByEmail(email string) (*models.User, error) {
	db := database.Database.Conn
	user := new(models.User)
	if err := db.Find(user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func checkEmailDuplicate(user *models.User) bool {
	database.Database.Conn.Find(user, "email = ?", user.Email)
	return user.ID > 0
}

func validUser(id string, p string) bool {
	db := database.Database.Conn
	user := new(models.User)
	db.First(user, id)
	fmt.Println("아이디가 없으면???")
	fmt.Printf("%#v", user)
	return CheckPasswordHash(p, user.Password)

}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))
	return uid == n
}
