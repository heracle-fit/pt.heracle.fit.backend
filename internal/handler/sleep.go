package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type SleepHandler struct {
	svc *service.SleepService
}

func NewSleepHandler(svc *service.SleepService) *SleepHandler {
	return &SleepHandler{svc: svc}
}

// AddSleepData godoc
// @Summary      Add sleep data
// @Description  Adds a new sleep data entry for the authenticated user (keeps last 7 entries)
// @Tags         Sleep
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body map[string]interface{} true "Sleep data entry"
// @Success      200 {object} map[string]interface{} "Updated sleep record"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /sleep/ [post]
func (h *SleepHandler) AddSleepData(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var newEntry map[string]interface{}
	if err := c.BodyParser(&newEntry); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.AddSleepData(context.Background(), userID, newEntry)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetSleepData godoc
// @Summary      Get sleep data
// @Description  Returns the last 7 sleep data entries for the authenticated user
// @Tags         Sleep
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Sleep data"
// @Router       /sleep/ [get]
func (h *SleepHandler) GetSleepData(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetSleepData(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetAIInsight godoc
// @Summary      Get AI sleep insight
// @Description  Returns an AI-generated insight based on recent sleep data (cached daily)
// @Tags         Sleep
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "AI insight"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /sleep/insight [get]
func (h *SleepHandler) GetAIInsight(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetAIInsight(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}
