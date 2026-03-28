package handler

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type SplitHandler struct {
	svc *service.SplitService
}

func NewSplitHandler(svc *service.SplitService) *SplitHandler {
	return &SplitHandler{svc: svc}
}

// GetMySplit godoc
// @Summary      Get my workout split
// @Description  Returns the workout split assigned to the authenticated user
// @Tags         Split
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Workout split"
// @Failure      404 {object} map[string]interface{} "Split not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /split/ [get]
func (h *SplitHandler) GetMySplit(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetMySplit(context.Background(), userID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// TrainerGetClientSplit godoc
// @Summary      Trainer: get client workout split
// @Description  Returns the workout split for a trainer's assigned client
// @Tags         Split
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Success      200 {object} map[string]interface{} "Client workout split"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      404 {object} map[string]interface{} "Split not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /split/trainer/{clientId} [get]
func (h *SplitHandler) TrainerGetClientSplit(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	result, err := h.svc.TrainerGetClientSplit(context.Background(), trainerUserID, clientID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// UpsertClientSplit godoc
// @Summary      Trainer: create or update client workout split
// @Description  Creates or updates the workout split for a trainer's assigned client
// @Tags         Split
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        body body object true "Split data" example({"splitData":{}})
// @Success      200 {object} map[string]interface{} "Upserted workout split"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /split/trainer/{clientId} [put]
func (h *SplitHandler) UpsertClientSplit(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")

	var body struct {
		SplitData json.RawMessage `json:"splitData"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}

	result, err := h.svc.UpsertClientSplit(context.Background(), trainerUserID, clientID, body.SplitData)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}
