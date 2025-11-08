package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// StatsHandler handles statistics-related HTTP requests
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler creates a new stats handler instance
func NewStatsHandler() (*StatsHandler, error) {
	statsService, err := service.NewStatsService()
	if err != nil {
		return nil, err
	}

	return &StatsHandler{
		statsService: statsService,
	}, nil
}

// GetCountries handles GET /api/stats/countries
// @Summary Get country wage statistics
// @Description Retrieve all countries with average wages and work time calculations
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/stats/countries [get]
func (h *StatsHandler) GetCountries(c *fiber.Ctx) error {
	countries := h.statsService.GetCountries()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"countries": countries,
			"count":     len(countries),
		},
	})
}

// GetCountryByCode handles GET /api/stats/countries/:code
// @Summary Get specific country statistics
// @Description Retrieve statistics for a specific country by its code
// @Tags stats
// @Accept json
// @Produce json
// @Param code path string true "Country Code (e.g., US, GB, RU)"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/stats/countries/{code} [get]
func (h *StatsHandler) GetCountryByCode(c *fiber.Ctx) error {
	code := c.Params("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "country code is required",
		})
	}

	country, err := h.statsService.GetCountryByCode(code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "country not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    country,
	})
}
