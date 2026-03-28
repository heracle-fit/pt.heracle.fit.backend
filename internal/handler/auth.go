package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/heracle/pt.heracle.fit.go/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// AuthenticateGoogleToken godoc
// @Summary      Authenticate with Google ID token
// @Description  Verifies a Firebase ID token and creates or retrieves the user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body object true "Google auth payload" example({"idToken":"...","accessToken":"..."})
// @Success      200 {object} map[string]interface{} "User object and JWT token"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      401 {object} map[string]interface{} "Invalid token"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /auth/google/token [post]
func (h *AuthHandler) AuthenticateGoogleToken(c *fiber.Ctx) error {
	var body struct {
		IDToken     string `json:"idToken"`
		AccessToken string `json:"accessToken"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}

	result, err := h.svc.AuthenticateGoogleToken(c.Context(), body.IDToken, body.AccessToken)
	log.Println("err: ", err)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"statusCode": 401, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"user": result.User, "token": result.Token})
}

// AdminLogin godoc
// @Summary      Admin login
// @Description  Authenticates admin with username and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body object true "Admin credentials" example({"username":"admin","password":"..."})
// @Success      200 {object} map[string]interface{} "Admin user and JWT token"
// @Failure      400 {object} map[string]interface{} "Invalid request body"
// @Failure      401 {object} map[string]interface{} "Invalid credentials"
// @Router       /auth/admin/login [post]
func (h *AuthHandler) AdminLogin(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "Invalid request body"})
	}

	result, err := h.svc.AdminLogin(c.Context(), body.Username, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"statusCode": 401, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"user": result.User, "token": result.Token})
}

// GetDevToken godoc
// @Summary      Get development token
// @Description  Creates or retrieves a user by email and returns a JWT (dev only)
// @Tags         Auth
// @Produce      json
// @Param        email query string false "Email address" default(sanjaysagar.main@gmail.com)
// @Success      200 {object} map[string]interface{} "User object and JWT token"
// @Failure      500 {object} map[string]interface{} "Server error"
// @Router       /auth/dev/token [get]
func (h *AuthHandler) GetDevToken(c *fiber.Ctx) error {
	email := c.Query("email", "sanjaysagar.main@gmail.com")
	result, err := h.svc.GetDevToken(c.Context(), email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}
	return c.JSON(fiber.Map{"user": result.User, "token": result.Token})
}
