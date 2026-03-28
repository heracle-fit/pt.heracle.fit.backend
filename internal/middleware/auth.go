package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Missing or invalid Authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Invalid token claims",
			})
		}

		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"statusCode": 401,
					"message":    "Token has expired",
				})
			}
		}

		c.Locals("userId", claims["sub"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])
		if username, ok := claims["username"].(string); ok {
			c.Locals("username", username)
		}

		return c.Next()
	}
}

func GenerateJWT(secret string, sub, email, role string, username ...string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   sub,
		"email": email,
		"role":  role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}
	if len(username) > 0 && username[0] != "" {
		claims["username"] = username[0]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
