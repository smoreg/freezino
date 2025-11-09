package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/model"
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

// StartWorkRequest represents the request body for starting work
type StartWorkRequest struct {
	JobType string `json:"job_type"`
}

// StartWork handles POST /api/work/start
// @Summary Start a work session
// @Description Start a new work session for the authenticated user
// @Tags work
// @Accept json
// @Produce json
// @Param request body StartWorkRequest true "Job type selection"
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

	// Parse request body
	var req StartWorkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	// Default to office if not specified
	if req.JobType == "" {
		req.JobType = "office"
	}

	result, err := h.workService.StartWork(userID, model.JobType(req.JobType))
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
		// Check for job requirement errors
		if err.Error() == "office_no_clothes" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "office_no_clothes",
			})
		}
		if err.Error() == "courier_no_uniform" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "courier_no_uniform",
			})
		}
		if err.Error() == "stunt_driver_no_car" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "stunt_driver_no_car",
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

// GetAvailableJobs handles GET /api/work/jobs
// @Summary Get available job types
// @Description Get list of available job types with requirements
// @Tags work
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/work/jobs [get]
func (h *WorkHandler) GetAvailableJobs(c *fiber.Ctx) error {
	jobs := []map[string]interface{}{
		{
			"type":        "office",
			"name":        "Office Worker",
			"base_reward": 500,
			"requires":    "clothing",
			"description": "Standard office job. Requires at least one clothing item.",
		},
		{
			"type":        "courier",
			"name":        "Courier",
			"base_reward": 500,
			"requires":    "courier_uniform",
			"bonus":       "own_car",
			"description": "Delivery work. Requires Courier Uniform. +$250 with own car.",
		},
		{
			"type":        "lab_rat",
			"name":        "Lab Test Subject",
			"base_reward": 500,
			"reward":      "random_mutation",
			"description": "Be a test subject for mad scientist. Receive random mutation.",
		},
		{
			"type":        "stunt_driver",
			"name":        "Stunt Driver",
			"base_reward": 1500,
			"requires":    "car",
			"penalty":     "car_broken",
			"description": "High-risk stunts. Requires car. Earn $1500 but car gets broken.",
		},
		{
			"type":        "drug_dealer",
			"name":        "Surprise Delivery",
			"base_reward": 2000,
			"penalty":     "jail_8years",
			"description": "Risky business. Earn $2000 but you'll be caught (8 year sentence, skippable).",
		},
		{
			"type":        "streamer",
			"name":        "Streamer",
			"base_reward": 0,
			"variable":    "lottery",
			"description": "Stream online. 70% = $0, 29% = $1, 1% = $10,000 + go viral (all future streams $10k).",
		},
		{
			"type":        "bottle_collector",
			"name":        "Bottle Collector",
			"base_reward": 100,
			"description": "Collect bottles and cans. Always $100. Available to everyone.",
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    jobs,
	})
}

// SkipJailTime handles POST /api/work/skip-jail
// @Summary Skip jail time
// @Description Skip the jail sentence time (instant release)
// @Tags work
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/work/skip-jail [post]
func (h *WorkHandler) SkipJailTime(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	if err := h.workService.SkipJailTime(userID); err != nil {
		if err.Error() == "not in jail" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "not in jail",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to skip jail time",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "jail time skipped successfully",
	})
}
