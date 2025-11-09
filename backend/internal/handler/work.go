package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// WorkHandler handles work-related HTTP requests
type WorkHandler struct {
	workService *service.WorkService
}

// NewWorkHandler creates a new work handler instance
func NewWorkHandler() *WorkHandler {
	return &WorkHandler{
		workService: service.NewWorkService(),
	}
}

// StartWork handles POST /api/work/start
// @Summary Start a work session
// @Description Start a new work session for the authenticated user
// @Tags work
// @Accept json
// @Produce json
// @Success 200 {object} service.StartWorkResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/work/start [post]
func (h *WorkHandler) StartWork(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	result, err := h.workService.StartWork(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		if err.Error() == "work session already in progress" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "work session already in progress",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to start work session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "work session started successfully",
	})
}

// GetStatus handles GET /api/work/status
// @Summary Get work session status
// @Description Get the current work session status for the authenticated user
// @Tags work
// @Accept json
// @Produce json
// @Success 200 {object} service.WorkStatusResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/work/status [get]
func (h *WorkHandler) GetStatus(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	status, err := h.workService.GetStatus(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get work status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    status,
	})
}

// CompleteWork handles POST /api/work/complete
// @Summary Complete work session
// @Description Complete the current work session and receive reward
// @Tags work
// @Accept json
// @Produce json
// @Success 200 {object} service.CompleteWorkResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/work/complete [post]
func (h *WorkHandler) CompleteWork(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	result, err := h.workService.CompleteWork(userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		if err.Error() == "no active work session" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   true,
				"message": "no active work session",
			})
		}
		// Check for "work not completed yet" error
		if len(err.Error()) > 23 && err.Error()[:23] == "work not completed yet," {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to complete work session",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "work session completed successfully",
	})
}

// GetHistory handles GET /api/work/history
// @Summary Get work session history
// @Description Get the work session history for the authenticated user
// @Tags work
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of sessions" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} service.WorkHistoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/work/history [get]
func (h *WorkHandler) GetHistory(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

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

	history, err := h.workService.GetHistory(uint(userID), limit, offset)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get work history",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    history,
	})
}
