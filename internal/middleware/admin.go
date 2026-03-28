package middleware

import "github.com/gofiber/fiber/v2"

func AdminGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": 403,
				"message":    "Admin privileges required",
			})
		}
		return c.Next()
	}
}
