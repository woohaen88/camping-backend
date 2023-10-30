package middleware

import (
	"camping-backend/config"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

//JWT는 JWT(JSON Web Token) 인증 미들웨어를 Return
//유효한 토큰의 경우 Ctx.Locals
//	1. 사용자를 설정
//	2. 핸들러를 호출
//잘못된 토큰의 경우 "401 - Unauthorized" 오류가 Return
//토큰이 누락된 경우 "400 - 잘못된 요청" 오류가 Return

func JwtMiddleWare() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
	})
}
