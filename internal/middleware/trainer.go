package middleware

import "github.com/gofiber/fiber/v2"

func TrainerGuard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if role != "trainer" && role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"statusCode": 403,
				"message":    "Trainer privileges required",
			})
		}
		return c.Next()
	}
}
