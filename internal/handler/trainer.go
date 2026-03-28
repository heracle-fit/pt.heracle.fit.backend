package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type TrainerHandler struct {
	svc *service.TrainerService
}

func NewTrainerHandler(svc *service.TrainerService) *TrainerHandler {
	return &TrainerHandler{svc: svc}
}

// GetClients godoc
// @Summary      Get trainer's clients
// @Description  Returns all clients assigned to the trainer with daily progress
// @Tags         Trainer
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} map[string]interface{} "Client list with progress"
// @Failure      403 {object} map[string]interface{} "Trainer not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /trainer/clients [get]
func (h *TrainerHandler) GetClients(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	result, err := h.svc.GetClients(context.Background(), trainerUserID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// GetClientDetails godoc
// @Summary      Get client details
// @Description  Returns detailed profile, metrics, and progress for a specific client
// @Tags         Trainer
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Success      200 {object} map[string]interface{} "Client details"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      404 {object} map[string]interface{} "Client not found"
// @Router       /trainer/client/{clientId} [get]
func (h *TrainerHandler) GetClientDetails(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	result, err := h.svc.GetClientDetails(context.Background(), trainerUserID, clientID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// AddClient godoc
// @Summary      Add a client
// @Description  Assigns a user as a client to the trainer by email
// @Tags         Trainer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object true "Client email" example({"email":"client@example.com"})
// @Success      200 {object} map[string]interface{} "Added client with progress"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      403 {object} map[string]interface{} "Trainer not found"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      409 {object} map[string]interface{} "Already assigned"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /trainer/clients/add [post]
func (h *TrainerHandler) AddClient(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	var body struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.AddClient(context.Background(), trainerUserID, body.Email)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// RemoveClient godoc
// @Summary      Remove a client
// @Description  Removes a client assignment from the trainer
// @Tags         Trainer
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Success      200 "Removed"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /trainer/clients/remove/{clientId} [delete]
func (h *TrainerHandler) RemoveClient(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	if err := h.svc.RemoveClient(context.Background(), trainerUserID, clientID); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(200)
}

// AdminAddTrainer godoc
// @Summary      Admin: add a trainer
// @Description  Promotes an existing user to trainer role (admin only)
// @Tags         Trainer
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body object true "Trainer details" example({"email":"trainer@example.com","specialization":"strength","experience":5})
// @Success      200 {object} map[string]interface{} "Created trainer with user info"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      409 {object} map[string]interface{} "Already a trainer"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /trainer/admin/add [post]
func (h *TrainerHandler) AdminAddTrainer(c *fiber.Ctx) error {
	var body struct {
		Email          string  `json:"email"`
		Specialization *string `json:"specialization"`
		Experience     *int    `json:"experience"`
	}
	c.BodyParser(&body)

	result, err := h.svc.AdminAddTrainer(context.Background(), body.Email, body.Specialization, body.Experience)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}
