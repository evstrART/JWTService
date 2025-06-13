package middleware

import (
	"JWTService/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authService *service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "no authorization token",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token type: expected access token",
			})
		}

		jti, ok := claims["jti"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid token format"})
		}

		redisKey := "revoked:" + jti
		vals, err := authService.RedisClient.HGetAll(c.Context(), redisKey).Result()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
		if len(vals) == 0 {
			return c.Status(401).JSON(fiber.Map{"error": "token revoked or not found"})
		}
		if vals["revoked"] == "1" {
			return c.Status(401).JSON(fiber.Map{"error": "token revoked"})
		}

		c.Context().SetUserValue("access_token", parts[1])
		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}
