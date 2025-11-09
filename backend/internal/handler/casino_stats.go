package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// CasinoStatsHandler handles casino statistics HTTP requests
type CasinoStatsHandler struct {
	casinoStatsService *service.CasinoStatsService
}

// NewCasinoStatsHandler creates a new casino stats handler instance
func NewCasinoStatsHandler() *CasinoStatsHandler {
	return &CasinoStatsHandler{
		casinoStatsService: service.NewCasinoStatsService(),
	}
}

// GetCasinoStats handles GET /api/casino/stats
// @Summary Get casino statistics
// @Description Retrieve overall casino statistics including total bets, wins, house edge, player profitability, etc.
// @Tags casino
// @Accept json
// @Produce json
// @Success 200 {object} service.CasinoStatsResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/casino/stats [get]
func (h *CasinoStatsHandler) GetCasinoStats(c *fiber.Ctx) error {
	// Get casino stats from service
	stats, err := h.casinoStatsService.GetCasinoStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get casino statistics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    stats,
	})
}
