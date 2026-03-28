package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type WorkoutHandler struct {
	svc *service.WorkoutService
}

func NewWorkoutHandler(svc *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{svc: svc}
}

// GetTodayWorkout godoc
// @Summary      Get today's workout suggestion
// @Description  Returns a workout suggestion based on user preferences
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Workout suggestion"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/today [get]
func (h *WorkoutHandler) GetTodayWorkout(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetTodayWorkout(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetStaticSessions godoc
// @Summary      Get static workout sessions
// @Description  Returns hardcoded sample workout sessions
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} map[string]interface{} "Static sessions"
// @Router       /workout/sessions [get]
func (h *WorkoutHandler) GetStaticSessions(c *fiber.Ctx) error {
	return c.JSON(h.svc.GetStaticSessions())
}

// GetExercises godoc
// @Summary      Get all exercises
// @Description  Returns the full exercise library with image URLs
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} map[string]interface{} "Exercise list"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/exercises [get]
func (h *WorkoutHandler) GetExercises(c *fiber.Ctx) error {
	result, err := h.svc.GetExercises(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetWorkoutPreferences godoc
// @Summary      Get workout preferences
// @Description  Returns workout preferences for the authenticated user
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Workout preferences"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/preferences [get]
func (h *WorkoutHandler) GetWorkoutPreferences(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetWorkoutPreferences(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// SaveWorkoutPreferences godoc
// @Summary      Save workout preferences
// @Description  Saves or updates workout preferences for the authenticated user
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body map[string]interface{} true "Workout preferences"
// @Success      200 {object} map[string]interface{} "Updated preferences"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/preferences [post]
func (h *WorkoutHandler) SaveWorkoutPreferences(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.SaveWorkoutPreferences(context.Background(), userID, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// ── Session CRUD ────────────────────────────────────────────────────────────────

// CreateSession godoc
// @Summary      Create a workout session
// @Description  Creates a new custom workout session for the authenticated user
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body service.CreateSessionRequest true "Session data"
// @Success      200 {object} map[string]interface{} "Created session"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/session [post]
func (h *WorkoutHandler) CreateSession(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var req service.CreateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.CreateSession(context.Background(), userID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetSession godoc
// @Summary      Get a workout session
// @Description  Returns a specific workout session by ID
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Session ID"
// @Success      200 {object} map[string]interface{} "Session details"
// @Failure      404 {object} map[string]interface{} "Session not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/session/{id} [get]
func (h *WorkoutHandler) GetSession(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	id, _ := strconv.Atoi(c.Params("id"))
	result, err := h.svc.GetSession(context.Background(), userID, id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// GetUserSessions godoc
// @Summary      Get user's workout sessions
// @Description  Returns all custom workout sessions for the authenticated user
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} map[string]interface{} "User sessions"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/my-sessions [get]
func (h *WorkoutHandler) GetUserSessions(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetUserSessions(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// UpdateSession godoc
// @Summary      Update a workout session
// @Description  Updates an existing workout session by ID
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Session ID"
// @Param        body body map[string]interface{} true "Fields to update"
// @Success      200 {object} map[string]interface{} "Updated session"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      404 {object} map[string]interface{} "Session not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/session/{id} [patch]
func (h *WorkoutHandler) UpdateSession(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	id, _ := strconv.Atoi(c.Params("id"))
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.UpdateSession(context.Background(), userID, id, body)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// DeleteSession godoc
// @Summary      Delete a workout session
// @Description  Deletes a workout session by ID
// @Tags         Workout
// @Security     BearerAuth
// @Param        id path int true "Session ID"
// @Success      200 "Deleted"
// @Failure      404 {object} map[string]interface{} "Session not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/session/{id} [delete]
func (h *WorkoutHandler) DeleteSession(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.svc.DeleteSession(context.Background(), userID, id); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(200)
}

// ── Trainer Session operations ──────────────────────────────────────────────────

// TrainerUpdateSession godoc
// @Summary      Trainer: update client session
// @Description  Allows a trainer to update a workout session for their assigned client
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        sessionId path int true "Session ID"
// @Param        body body map[string]interface{} true "Fields to update"
// @Success      200 {object} map[string]interface{} "Updated session"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      404 {object} map[string]interface{} "Session not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/trainer/session/{clientId}/{sessionId} [patch]
func (h *WorkoutHandler) TrainerUpdateSession(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	sessionID, _ := strconv.Atoi(c.Params("sessionId"))

	var body map[string]interface{}
	c.BodyParser(&body)

	result, err := h.svc.TrainerUpdateSession(context.Background(), trainerUserID, clientID, sessionID, body)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// TrainerCreateSession godoc
// @Summary      Trainer: create client session
// @Description  Allows a trainer to create a workout session for their assigned client
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        body body service.CreateSessionRequest true "Session data"
// @Success      200 {object} map[string]interface{} "Created session"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/trainer/session/{clientId} [post]
func (h *WorkoutHandler) TrainerCreateSession(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")

	var req service.CreateSessionRequest
	c.BodyParser(&req)

	result, err := h.svc.TrainerCreateSession(context.Background(), trainerUserID, clientID, req)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// TrainerGetSessions godoc
// @Summary      Trainer: get client sessions
// @Description  Returns all workout sessions for a trainer's assigned client
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Success      200 {array} map[string]interface{} "Client sessions"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/trainer/sessions/{clientId} [get]
func (h *WorkoutHandler) TrainerGetSessions(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")

	result, err := h.svc.TrainerGetSessions(context.Background(), trainerUserID, clientID)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// TrainerDeleteSession godoc
// @Summary      Trainer: delete client session
// @Description  Allows a trainer to delete a workout session for their assigned client
// @Tags         Workout
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        sessionId path int true "Session ID"
// @Success      200 "Deleted"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      404 {object} map[string]interface{} "Session not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/trainer/session/{clientId}/{sessionId} [delete]
func (h *WorkoutHandler) TrainerDeleteSession(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	sessionID, _ := strconv.Atoi(c.Params("sessionId"))

	if err := h.svc.TrainerDeleteSession(context.Background(), trainerUserID, clientID, sessionID); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(200)
}

// ── WorkoutLog CRUD ─────────────────────────────────────────────────────────────

// CreateWorkoutLog godoc
// @Summary      Create a workout log
// @Description  Logs a completed workout for the authenticated user
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body service.CreateWorkoutLogRequest true "Workout log data"
// @Success      200 {object} map[string]interface{} "Created workout log"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/log [post]
func (h *WorkoutHandler) CreateWorkoutLog(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var req service.CreateWorkoutLogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.CreateWorkoutLog(context.Background(), userID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetWorkoutLog godoc
// @Summary      Get a workout log
// @Description  Returns a specific workout log by ID
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Workout log ID"
// @Success      200 {object} map[string]interface{} "Workout log"
// @Failure      404 {object} map[string]interface{} "Log not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/log/{id} [get]
func (h *WorkoutHandler) GetWorkoutLog(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	id, _ := strconv.Atoi(c.Params("id"))
	result, err := h.svc.GetWorkoutLog(context.Background(), userID, id)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// GetWorkoutLogs godoc
// @Summary      Get all workout logs
// @Description  Returns all workout logs for the authenticated user
// @Tags         Workout
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} map[string]interface{} "Workout logs"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/logs [get]
func (h *WorkoutHandler) GetWorkoutLogs(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetWorkoutLogs(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// UpdateWorkoutLog godoc
// @Summary      Update a workout log
// @Description  Updates an existing workout log by ID
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Workout log ID"
// @Param        body body map[string]interface{} true "Fields to update"
// @Success      200 {object} map[string]interface{} "Updated workout log"
// @Failure      404 {object} map[string]interface{} "Log not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/log/{id} [patch]
func (h *WorkoutHandler) UpdateWorkoutLog(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	id, _ := strconv.Atoi(c.Params("id"))
	var body map[string]interface{}
	c.BodyParser(&body)
	result, err := h.svc.UpdateWorkoutLog(context.Background(), userID, id, body)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// DeleteWorkoutLog godoc
// @Summary      Delete a workout log
// @Description  Deletes a workout log by ID
// @Tags         Workout
// @Security     BearerAuth
// @Param        id path int true "Workout log ID"
// @Success      200 "Deleted"
// @Failure      404 {object} map[string]interface{} "Log not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/log/{id} [delete]
func (h *WorkoutHandler) DeleteWorkoutLog(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.svc.DeleteWorkoutLog(context.Background(), userID, id); err != nil {
		return handleServiceError(c, err)
	}
	return c.SendStatus(200)
}

// TrainerAddLogReview godoc
// @Summary      Trainer: add review to workout log
// @Description  Allows a trainer to add a review to their client's workout log
// @Tags         Workout
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        logId path int true "Workout log ID"
// @Param        body body object true "Review data" example({"review":"Good form"})
// @Success      200 {object} map[string]interface{} "Updated workout log with review"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      404 {object} map[string]interface{} "Log not found"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /workout/trainer/log-review/{logId} [patch]
func (h *WorkoutHandler) TrainerAddLogReview(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	logID, _ := strconv.Atoi(c.Params("logId"))

	var body struct {
		Review string `json:"review"`
	}
	c.BodyParser(&body)

	result, err := h.svc.TrainerAddLogReview(context.Background(), trainerUserID, logID, body.Review)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}
