package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
)

// DevHandler handles development/testing endpoints
type DevHandler struct{}

// NewDevHandler creates a new dev handler
func NewDevHandler() *DevHandler {
	return &DevHandler{}
}

// AddMoney adds money to a user (dev/testing only)
func (h *DevHandler) AddMoney(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get amount from body or use default
	type AddMoneyRequest struct {
		Amount float64 `json:"amount"`
	}

	req := AddMoneyRequest{Amount: 1000.00} // Default amount
	_ = c.BodyParser(&req)

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "amount must be positive",
		})
	}

	// Get database
	db := database.GetDB()

	// Find user
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "user not found",
		})
	}

	// Add money
	user.Balance += req.Amount

	// Save
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to update balance",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":     user.ID,
			"new_balance": user.Balance,
			"added":       req.Amount,
		},
	})
}

// ResetBalance resets user balance to 1000 (dev/testing only)
func (h *DevHandler) ResetBalance(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get database
	db := database.GetDB()

	// Find user
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "user not found",
		})
	}

	// Reset balance
	user.Balance = 1000.00

	// Save
	if err := db.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to reset balance",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id":  user.ID,
			"balance": user.Balance,
			"message": "Balance reset to $1000",
		},
	})
}

// SeedDatabase runs database seeding (dev/testing only)
func (h *DevHandler) SeedDatabase(c *fiber.Ctx) error {
	// Run seed
	if err := database.Seed(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Database seeded successfully (test users and items created)",
	})
}
