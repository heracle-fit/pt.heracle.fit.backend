package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Returns the authenticated user's profile with full details
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "User profile"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	profile, err := h.svc.GetProfile(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	if profile == nil {
		return c.Status(404).JSON(fiber.Map{"statusCode": 404, "message": "User not found"})
	}
	return c.JSON(profile)
}

// GetBodyMetrics godoc
// @Summary      Get body metrics
// @Description  Returns the current body metrics for the authenticated user
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Body metrics"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/body-metrics [get]
func (h *UserHandler) GetBodyMetrics(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	metrics, err := h.svc.GetBodyMetrics(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(metrics)
}

// GetOnboardingStatus godoc
// @Summary      Get onboarding status
// @Description  Returns the onboarding completion status for the authenticated user
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Onboarding status"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/onboarding-status [get]
func (h *UserHandler) GetOnboardingStatus(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	status, err := h.svc.GetOnboardingStatus(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(status)
}

// TestCalendar godoc
// @Summary      Test Google Calendar
// @Description  Fetches the Google Calendar list for the authenticated user using their stored access token
// @Tags         User
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Google Calendar list"
// @Failure      403 {object} map[string]interface{} "Token expired"
// @Failure      404 {object} map[string]interface{} "User or token not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/caltest [get]
func (h *UserHandler) TestCalendar(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	body, status, err := h.svc.GetCalendarDetails(context.Background(), userID)
	if err != nil {
		return c.Status(status).JSON(fiber.Map{"statusCode": status, "message": err.Error()})
	}
	c.Set("Content-Type", "application/json")
	return c.Send(body)
}

// SaveBodyMetrics godoc
// @Summary      Save body metrics
// @Description  Saves or updates body metrics for the authenticated user. Automatically calculates BMI, maintenance calories, and macro targets.
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body map[string]interface{} true "Body metrics data"
// @Success      200 {object} map[string]interface{} "Updated body metrics"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/body-metrics [post]
func (h *UserHandler) SaveBodyMetrics(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}

	result, err := h.svc.SaveBodyMetrics(context.Background(), userID, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// SaveTargets godoc
// @Summary      Save nutritional targets
// @Description  Saves or updates nutritional targets (calories, macros) for the authenticated user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body map[string]interface{} true "Target data"
// @Success      200 {object} map[string]interface{} "Updated targets"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/targets [post]
func (h *UserHandler) SaveTargets(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.SaveTargets(context.Background(), userID, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// TrainerSaveBodyMetrics godoc
// @Summary      Trainer: save client body metrics
// @Description  Allows a trainer to save or update body metrics for their assigned client
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        body body map[string]interface{} true "Body metrics data"
// @Success      200 {object} map[string]interface{} "Updated body metrics"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /user/trainer/body-metrics/{clientId} [patch]
func (h *UserHandler) TrainerSaveBodyMetrics(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}

	result, err := h.svc.TrainerSaveBodyMetrics(context.Background(), trainerUserID, clientID, body)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}
