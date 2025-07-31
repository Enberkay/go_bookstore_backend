package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims := user.(map[string]interface{})
		if claims["role"] != role {
			return c.Status(403).JSON(fiber.Map{"error": "Forbidden: insufficient role"})
		}

		return c.Next()
	}
}
