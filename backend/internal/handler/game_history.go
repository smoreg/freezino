package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// GameHistoryHandler handles game history HTTP requests
type GameHistoryHandler struct {
	gameHistoryService *service.GameHistoryService
}

// NewGameHistoryHandler creates a new game history handler instance
func NewGameHistoryHandler() *GameHistoryHandler {
	return &GameHistoryHandler{
		gameHistoryService: service.NewGameHistoryService(),
	}
}

// GetHistory handles GET /api/games/history
// @Summary Get game history
// @Description Retrieve paginated game history with optional filters
// @Tags games
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param game query string false "Game type filter (roulette, slots, blackjack, etc.)"
// @Param limit query int false "Limit number of records" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} service.GameHistoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/history [get]
func (h *GameHistoryHandler) GetHistory(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	// Get optional game type filter
	gameType := c.Query("game", "")

	// Parse pagination parameters
	limit := 50 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // default offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get history from service
	history, err := h.gameHistoryService.GetHistory(userID, gameType, limit, offset)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get game history",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    history,
	})
}

// GetStats handles GET /api/games/stats
// @Summary Get game statistics
// @Description Retrieve overall game statistics for a user
// @Tags games
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {object} service.GameStatsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/stats [get]
func (h *GameHistoryHandler) GetStats(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	// Get stats from service
	stats, err := h.gameHistoryService.GetStats(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get game statistics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}
