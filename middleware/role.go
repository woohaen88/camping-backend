package middleware

import (
	"camping-backend/enums"
	"github.com/gofiber/fiber/v2"
)

func AssignRole(role enums.Role) fiber.Handler {

	if err := role.Check(); err != nil {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c *fiber.Ctx) error {
		authUser, err := GetAuthUser(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "credentials is invalid",
				"data":    err.Error(),
			})
		}

		if authUser.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "user role is not enough permission",
				"data":    nil,
			})
		}
		return nil
	}
}
