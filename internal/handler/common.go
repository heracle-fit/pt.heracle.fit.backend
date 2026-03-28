package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

// handleServiceError maps a service.AppError to the appropriate HTTP response.
func handleServiceError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*service.AppError); ok {
		return c.Status(appErr.Status).JSON(fiber.Map{"statusCode": appErr.Status, "message": appErr.Message})
	}
	return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
}
