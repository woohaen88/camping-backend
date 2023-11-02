package errors

import "github.com/gofiber/fiber/v2"

func ErrorHandler(c *fiber.Ctx, status int, err error, message ...string) error {
	m := "error"
	if len(message) > 0 {
		m = message[0]
	}

	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": m,
		"data":    err.Error(),
	})

}
