package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// SlotsHandler handles slots game HTTP requests
type SlotsHandler struct {
	slotsService *service.SlotsService
}

// NewSlotsHandler creates a new slots handler instance
func NewSlotsHandler() *SlotsHandler {
	return &SlotsHandler{
		slotsService: service.NewSlotsService(),
	}
}

// Spin handles POST /api/games/slots/spin
// @Summary Spin the slot machine
// @Description Spin the slot machine with a specified bet
// @Tags games
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param bet body number true "Bet amount"
// @Success 200 {object} service.SpinResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/games/slots/spin [post]
func (h *SlotsHandler) Spin(c *fiber.Ctx) error {
	// Get user_id from query parameter
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
	var reqBody struct {
		Bet float64 `json:"bet"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	if reqBody.Bet <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "bet must be greater than 0",
		})
	}

	// Perform spin
	spinReq := &service.SpinRequest{
		UserID: uint(userID),
		Bet:    reqBody.Bet,
	}

	result, err := h.slotsService.Spin(spinReq)
	if err != nil {
		// Check for specific error types
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}

		// Check for insufficient balance error
		if len(err.Error()) > 20 && err.Error()[:20] == "insufficient balance" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to spin slots",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "spin successful",
	})
}

// GetPayoutTable handles GET /api/games/slots/payouts
// @Summary Get payout table
// @Description Get the payout table for the slot machine
// @Tags games
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/games/slots/payouts [get]
func (h *SlotsHandler) GetPayoutTable(c *fiber.Ctx) error {
	payouts := h.slotsService.GetPayoutTable()
	symbols := h.slotsService.GetSymbols()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"payouts": payouts,
			"symbols": symbols,
		},
	})
}
