package middleware

import (
	"camping-backend/config"
	"camping-backend/database"
	"camping-backend/models"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

//JWT는 JWT(JSON Web Token) 인증 미들웨어를 Return
//유효한 토큰의 경우 Ctx.Locals
//	1. 사용자를 설정
//	2. 핸들러를 호출
//잘못된 토큰의 경우 "401 - Unauthorized" 오류가 Return
//토큰이 누락된 경우 "400 - 잘못된 요청" 오류가 Return

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
	})
}

func GetAuthUser(c *fiber.Ctx) (*models.User, error) {

	jwtUser := c.Locals("user").(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	userId := claims["userId"]

	db := database.Database.Conn
	user := new(models.User)

	if err := db.First(user, "id = ?", userId).Error; err != nil {

		return nil, err
	}

	return user, nil

}
