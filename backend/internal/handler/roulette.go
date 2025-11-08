package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// RouletteHandler handles roulette game HTTP requests
type RouletteHandler struct {
	rouletteService *service.RouletteService
}

// NewRouletteHandler creates a new roulette handler instance
func NewRouletteHandler() *RouletteHandler {
	return &RouletteHandler{
		rouletteService: service.NewRouletteService(),
	}
}

// PlaceBet handles POST /api/games/roulette/bet
// @Summary Place a roulette bet
// @Description Place one or more bets on the roulette table and spin
// @Tags roulette
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param request body service.PlaceBetRequest true "Bet request"
// @Success 200 {object} service.PlaceBetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/roulette/bet [post]
func (h *RouletteHandler) PlaceBet(c *fiber.Ctx) error {
	// Get user ID from query parameter (in production, use JWT)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user_id",
		})
	}

	// Parse request body
	var req service.PlaceBetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	// Override user_id from query parameter
	req.UserID = uint(userID)

	// Validate bets
	if len(req.Bets) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "at least one bet is required",
		})
	}

	// Place bet
	result, err := h.rouletteService.PlaceBet(req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		case "insufficient balance":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "insufficient balance",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "failed to place bet: " + err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "bet placed successfully",
	})
}

// GetHistory handles GET /api/games/roulette/history
// @Summary Get roulette game history
// @Description Retrieve recent roulette game history for a user
// @Tags roulette
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param limit query int false "Limit number of results" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/roulette/history [get]
func (h *RouletteHandler) GetHistory(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user_id",
		})
	}

	// Parse limit parameter
	limit := 10 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	history, err := h.rouletteService.GetHistory(uint(userID), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get history",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"history": history,
			"count":   len(history),
		},
	})
}

// GetRecentNumbers handles GET /api/games/roulette/recent
// @Summary Get recent winning numbers
// @Description Retrieve recent winning numbers across all players
// @Tags roulette
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/roulette/recent [get]
func (h *RouletteHandler) GetRecentNumbers(c *fiber.Ctx) error {
	// Parse limit parameter
	limit := 20 // default limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	numbers, err := h.rouletteService.GetRecentNumbers(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get recent numbers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"numbers": numbers,
			"count":   len(numbers),
		},
	})
}
