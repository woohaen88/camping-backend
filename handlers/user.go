package handlers

import (
	"camping-backend/database"
	"camping-backend/models"
	"camping-backend/serializers"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func validUser(id string, p string) bool {
	db := database.DB
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
	user.PaswordHash()

	database.DB.Create(user)
	responseUser := serializers.UserSerializer(*user)
	return c.Status(200).JSON(responseUser)

}

func checkEmailDuplicate(user *models.User) bool{
	database.DB.Find(user, "email = ?", user.Email)	
	return user.ID>0
}

