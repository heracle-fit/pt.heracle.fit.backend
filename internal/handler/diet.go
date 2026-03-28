package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type DietHandler struct {
	svc *service.DietService
}

func NewDietHandler(svc *service.DietService) *DietHandler {
	return &DietHandler{svc: svc}
}

// GetStatus godoc
// @Summary      Get daily nutritional status
// @Description  Returns target vs consumed macros for a given date
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Param        date query string false "Date in YYYY-MM-DD format" default(today)
// @Success      200 {object} map[string]interface{} "Nutritional status with targets and consumed"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/status [get]
func (h *DietHandler) GetStatus(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	date := c.Query("date", time.Now().Format("2006-01-02"))
	result, err := h.svc.GetStatus(context.Background(), userID, date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetTodayDiet godoc
// @Summary      Get today's diet suggestion
// @Description  Returns the AI-generated diet suggestion for today, generating one if none exists
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Diet suggestion"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/today [get]
func (h *DietHandler) GetTodayDiet(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetTodayDiet(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// SearchFood godoc
// @Summary      Search food items
// @Description  Searches the food database by name
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Param        query query string true "Search query"
// @Success      200 {array} map[string]interface{} "Matching food items"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/food/search [get]
func (h *DietHandler) SearchFood(c *fiber.Ctx) error {
	query := c.Query("query")
	result, err := h.svc.SearchFood(context.Background(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetMealsByDate godoc
// @Summary      Get meals by date
// @Description  Returns all logged meals for the authenticated user on a specific date
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Param        date query string false "Date in YYYY-MM-DD format" default(today)
// @Success      200 {array} map[string]interface{} "Meals for the date"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/meals [get]
func (h *DietHandler) GetMealsByDate(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	date := c.Query("date", time.Now().Format("2006-01-02"))
	result, err := h.svc.GetMealsByDate(context.Background(), userID, date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// GetDietPreferences godoc
// @Summary      Get diet preferences
// @Description  Returns diet preferences for the authenticated user
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{} "Diet preferences"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/preferences [get]
func (h *DietHandler) GetDietPreferences(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	result, err := h.svc.GetDietPreferences(context.Background(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// SaveDietPreferences godoc
// @Summary      Save diet preferences
// @Description  Saves or updates diet preferences for the authenticated user
// @Tags         Diet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body map[string]interface{} true "Diet preferences"
// @Success      200 {object} map[string]interface{} "Updated preferences"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/preferences [post]
func (h *DietHandler) SaveDietPreferences(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.SaveDietPreferences(context.Background(), userID, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// LogMeal godoc
// @Summary      Log a meal
// @Description  Logs a meal entry and triggers background AI diet suggestion regeneration
// @Tags         Diet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body service.LogMealRequest true "Meal data"
// @Success      200 {object} map[string]interface{} "Logged meal"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/meal [post]
func (h *DietHandler) LogMeal(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	var req service.LogMealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}
	result, err := h.svc.LogMeal(context.Background(), userID, req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(result)
}

// AnalyseFood godoc
// @Summary      Analyse food with AI
// @Description  Analyses food from an image and/or text description using AI
// @Tags         Diet
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        description formData string false "Text description of the food"
// @Param        image formData file false "Image of the food"
// @Success      200 {object} map[string]interface{} "Food analysis result"
// @Failure      400 {object} map[string]interface{} "Missing input or AI error"
// @Router       /diet/ai/food [post]
func (h *DietHandler) AnalyseFood(c *fiber.Ctx) error {
	description := c.FormValue("description")
	file, err := c.FormFile("image")

	var imageData []byte
	var mimeType string

	if err == nil && file != nil {
		f, err := file.Open()
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Failed to read image"})
		}
		defer f.Close()

		imageData = make([]byte, file.Size)
		f.Read(imageData)
		mimeType = file.Header.Get("Content-Type")
	}

	var desc *string
	if description != "" {
		desc = &description
	}

	result, svcErr := h.svc.AnalyseFood(context.Background(), desc, imageData, mimeType)
	if svcErr != nil {
		return handleServiceError(c, svcErr)
	}
	return c.JSON(result)
}

// TrainerUpdateTargets godoc
// @Summary      Trainer: update client targets
// @Description  Allows a trainer to update nutritional targets for their assigned client
// @Tags         Diet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        body body map[string]interface{} true "Target data"
// @Success      200 {object} map[string]interface{} "Updated targets"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/trainer/targets/{clientId} [patch]
func (h *DietHandler) TrainerUpdateTargets(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}

	result, err := h.svc.TrainerUpdateTargets(context.Background(), trainerUserID, clientID, body)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// TrainerGetMealsByDate godoc
// @Summary      Trainer: get client meals by date
// @Description  Returns all logged meals for a client on a specific date
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        date query string false "Date in YYYY-MM-DD format" default(today)
// @Success      200 {array} map[string]interface{} "Client meals"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/trainer/meals/{clientId} [get]
func (h *DietHandler) TrainerGetMealsByDate(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	date := c.Query("date", time.Now().Format("2006-01-02"))

	result, err := h.svc.TrainerGetMealsByDate(context.Background(), trainerUserID, clientID, date)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}

// TrainerGetStatus godoc
// @Summary      Trainer: get client nutritional status
// @Description  Returns target vs consumed macros for a client on a specific date
// @Tags         Diet
// @Produce      json
// @Security     BearerAuth
// @Param        clientId path string true "Client user ID"
// @Param        date query string false "Date in YYYY-MM-DD format" default(today)
// @Success      200 {object} map[string]interface{} "Client nutritional status"
// @Failure      403 {object} map[string]interface{} "Not authorized"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /diet/trainer/status/{clientId} [get]
func (h *DietHandler) TrainerGetStatus(c *fiber.Ctx) error {
	trainerUserID := c.Locals("userId").(string)
	clientID := c.Params("clientId")
	date := c.Query("date", time.Now().Format("2006-01-02"))

	result, err := h.svc.TrainerGetStatus(context.Background(), trainerUserID, clientID, date)
	if err != nil {
		return handleServiceError(c, err)
	}
	return c.JSON(result)
}
