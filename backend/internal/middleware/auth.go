package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/auth"
	"github.com/smoreg/freezino/backend/internal/config"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	jwtManager := auth.NewJWTManager(cfg)

	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			if errors.Is(err, auth.ErrExpiredToken) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "token has expired",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Check token type
		if claims.Type != auth.AccessToken {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token type",
			})
		}

		// Get user from database
		db := database.GetDB()
		var user model.User

		if err := db.First(&user, claims.UserID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "user not found",
			})
		}

		// Store user in context
		c.Locals("user", &user)
		c.Locals("userID", user.ID)

		return c.Next()
	}
}

// OptionalAuth is a middleware that adds user info if authenticated but doesn't require it
func OptionalAuth(cfg *config.Config) fiber.Handler {
	jwtManager := auth.NewJWTManager(cfg)

	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		token := parts[1]

		// Validate token
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			return c.Next()
		}

		// Check token type
		if claims.Type != auth.AccessToken {
			return c.Next()
		}

		// Get user from database
		db := database.GetDB()
		var user model.User

		if err := db.First(&user, claims.UserID).Error; err != nil {
			return c.Next()
		}

		// Store user in context
		c.Locals("user", &user)
		c.Locals("userID", user.ID)

		return c.Next()
	}
}
