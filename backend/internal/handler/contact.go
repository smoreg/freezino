package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
)

// ContactHandler handles contact-related HTTP requests
type ContactHandler struct{}

// NewContactHandler creates a new contact handler instance
func NewContactHandler() *ContactHandler {
	return &ContactHandler{}
}

// ContactRequest represents the request body for submitting a contact message
type ContactRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Email   string `json:"email" validate:"required,email,max=255"`
	Message string `json:"message" validate:"required,min=10,max=2000"`
}

// SubmitMessage handles POST /api/contact
// @Summary Submit a contact message
// @Description Submit a new contact message from the contact form
// @Tags contact
// @Accept json
// @Produce json
// @Param request body ContactRequest true "Contact Message"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/contact [post]
func (h *ContactHandler) SubmitMessage(c *fiber.Ctx) error {
	var req ContactRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" || len(req.Name) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Name must be at least 2 characters long",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Email is required",
		})
	}

	if req.Message == "" || len(req.Message) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Message must be at least 10 characters long",
		})
	}

	if len(req.Message) > 2000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Message must not exceed 2000 characters",
		})
	}

	// Create contact message
	contactMessage := model.ContactMessage{
		Name:    req.Name,
		Email:   req.Email,
		Message: req.Message,
		IsRead:  false,
	}

	// Save to database
	db := database.DB
	if err := db.Create(&contactMessage).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save contact message",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Your message has been sent successfully! We'll get back to you soon.",
		"data": fiber.Map{
			"id": contactMessage.ID,
		},
	})
}
