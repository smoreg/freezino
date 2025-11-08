package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/smoreg/freezino/backend/internal/auth"
	"github.com/smoreg/freezino/backend/internal/config"
	"github.com/smoreg/freezino/backend/internal/service"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
	jwtManager  *auth.JWTManager
	cfg         *config.Config
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtManager:  auth.NewJWTManager(cfg),
		cfg:         cfg,
	}
}

// generateTokens generates access and refresh JWT tokens using JWTManager
func (h *AuthHandler) generateTokens(userID uint, email string) (string, string, error) {
	return h.jwtManager.GenerateTokenPair(userID, email)
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req service.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Register user
	user, err := h.authService.Register(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Generate tokens
	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate tokens",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Authenticate user
	user, err := h.authService.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Generate tokens
	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate tokens",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate refresh token
	claims, err := h.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid or expired refresh token",
		})
	}

	// Check token type
	if claims.Type != auth.RefreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid token type",
		})
	}

	// Generate new access token
	accessToken, err := h.jwtManager.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to generate access token",
		})
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"access_token": accessToken,
	})
}
