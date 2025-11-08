package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HealthCheck handles health check requests
func HealthCheck(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "freezino-backend",
		Version:   "1.0.0",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
