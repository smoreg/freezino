package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/config"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
)

// Handler handles authentication requests
type Handler struct {
	config     *config.Config
	jwtManager *JWTManager
}

// NewHandler creates a new auth handler
func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		config:     cfg,
		jwtManager: NewJWTManager(cfg),
	}
}

// generateState generates a random state string for OAuth
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// GoogleLogin initiates Google OAuth flow
func (h *Handler) GoogleLogin(c *fiber.Ctx) error {
	oauthConfig := GoogleOAuthConfig(
		h.config.GoogleClientID,
		h.config.GoogleClientSecret,
		h.config.GoogleRedirectURL,
	)

	state := generateState()

	// Store state in session/cookie for verification in callback
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Lax",
	})

	url := oauthConfig.AuthCodeURL(state)
	return c.Redirect(url)
}

// GoogleCallback handles Google OAuth callback
func (h *Handler) GoogleCallback(c *fiber.Ctx) error {
	// Verify state
	state := c.Query("state")
	cookieState := c.Cookies("oauth_state")

	if state == "" || state != cookieState {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid state parameter",
		})
	}

	// Clear state cookie
	c.ClearCookie("oauth_state")

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing authorization code",
		})
	}

	oauthConfig := GoogleOAuthConfig(
		h.config.GoogleClientID,
		h.config.GoogleClientSecret,
		h.config.GoogleRedirectURL,
	)

	// Get user info from Google
	userInfo, err := GetGoogleUserInfo(c.Context(), oauthConfig, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get user info from Google",
		})
	}

	// Find or create user
	db := database.GetDB()
	var user model.User

	result := db.Where("google_id = ?", userInfo.ID).First(&user)
	if result.Error != nil {
		// User doesn't exist, create new one
		user = model.User{
			GoogleID: userInfo.ID,
			Email:    userInfo.Email,
			Name:     userInfo.Name,
			Avatar:   userInfo.Picture,
			Balance:  1000.00, // Initial balance
		}

		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to create user",
			})
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := h.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate tokens",
		})
	}

	// Set refresh token as HTTP-only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Lax",
		Path:     "/api/auth",
	})

	// Redirect to frontend with access token
	redirectURL := h.config.FrontendURL + "/auth/callback?token=" + accessToken
	return c.Redirect(redirectURL)
}

// GetMe returns the current authenticated user
func (h *Handler) GetMe(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	user := c.Locals("user").(*model.User)

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// Logout logs out the user by clearing tokens
func (h *Handler) Logout(c *fiber.Ctx) error {
	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   h.config.Environment == "production",
		SameSite: "Lax",
		Path:     "/api/auth",
	})

	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}

// RefreshToken generates a new access token using refresh token
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "refresh token not found",
		})
	}

	// Validate refresh token
	claims, err := h.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid refresh token",
		})
	}

	// Check token type
	if claims.Type != RefreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token type",
		})
	}

	// Generate new access token
	accessToken, err := h.jwtManager.GenerateAccessToken(claims.UserID, claims.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate access token",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}
