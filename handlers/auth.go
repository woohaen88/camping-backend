package handlers

import (
	"camping-backend/database"
	"camping-backend/models"
	"errors"
	"net/mail"
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


func Login (c *fiber.Ctx) error {
	type LoginInput struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type UserData struct {
		ID uint `json:"id"`
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	input := new(LoginInput)
	var userData UserData

	if err:= c.BodyParser(input); err!= nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "error", 
			"message": "Error on login request", 
			"data": err.Error(),
		})
	}

	email := input.Email
	password :=  input.Password
	userModel, err := new(models.User), *new(error)


	if isEmail(email) {
		userModel, err = getUserByEmail(email)
		} 
	
	if userModel == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", 
			"message": "User not found", 
			"data": err.Error(),
		})
	}

	userData = UserData{
		ID: userModel.ID,
		Username: userModel.Username,
		Email: userModel.Email,
		Password: userModel.Password,
	}

	if !CheckPasswordHash(password, userData.Password){
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error", 
			"message": "Invalid password", 
			"data": nil,
		})

	}
	// jwt
	token := jwt.New(jwt.SigningMethodHS256)
	claime := token.Claims.(jwt.MapClaims)
	claime["userid"] = userData.ID
	claime["exp"] = time.Now().Add(time.Hour*72).Unix()

	t, err := token.SignedString([]byte("mysecret")) // Todo


	
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success", 
		"message": "Success login", 
		"data": t,
	})
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
func getUserByEmail(email string) (*models.User, error) {
	db := database.DB
	user := new(models.User)
	if err := db.Find(user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}