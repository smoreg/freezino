package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// LoanHandler handles loan-related HTTP requests
type LoanHandler struct {
	loanService *service.LoanService
}

// NewLoanHandler creates a new loan handler instance
func NewLoanHandler() *LoanHandler {
	return &LoanHandler{
		loanService: service.NewLoanService(),
	}
}

// GetLoanSummary handles GET /api/loans/summary
// @Summary Get loan summary
// @Description Get aggregate loan information for the current user
// @Tags loans
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/loans/summary [get]
func (h *LoanHandler) GetLoanSummary(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	summary, err := h.loanService.GetLoanSummary(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get loan summary",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    summary,
	})
}

// GetUserLoans handles GET /api/loans
// @Summary Get user loans
// @Description Get all active loans for the current user
// @Tags loans
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/loans [get]
func (h *LoanHandler) GetUserLoans(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	loans, err := h.loanService.GetUserLoans(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get loans",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"loans": loans,
			"count": len(loans),
		},
	})
}

// TakeLoan handles POST /api/loans/take
// @Summary Take a loan
// @Description Apply for a new loan (friends, bank, or microcredit)
// @Tags loans
// @Accept json
// @Produce json
// @Param body body service.TakeLoanRequest true "Loan application"
// @Success 200 {object} service.TakeLoanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/loans/take [post]
func (h *LoanHandler) TakeLoan(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	// Parse request body
	var req service.TakeLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	// Set user ID from auth context
	req.UserID = userID

	// Process loan
	result, err := h.loanService.TakeLoan(req)
	if err != nil {
		errMsg := err.Error()

		// Handle specific error cases
		status := fiber.StatusBadRequest
		switch errMsg {
		case "friends_refused":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "friends_refused",
				"details": "Your friends are tired of lending you money",
			})
		case "friends_limit_exceeded":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "friends_limit_exceeded",
				"details": "You've already borrowed too much from friends",
			})
		case "collateral_required":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "collateral_required",
				"details": "Bank loans require collateral (car or house)",
			})
		case "collateral_not_found":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "collateral_not_found",
			})
		case "not_your_item":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "not_your_item",
			})
		case "item_already_collateral":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "item_already_collateral",
				"details": "This item is already used as collateral",
			})
		case "collateral_must_be_car_or_house":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "collateral_must_be_car_or_house",
				"details": "Only cars and houses can be used as collateral",
			})
		case "collateral_insufficient":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "collateral_insufficient",
				"details": "Collateral value is less than loan amount",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "failed to process loan",
				"details": errMsg,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// RepayLoan handles POST /api/loans/repay/:loanId
// @Summary Repay a loan
// @Description Make a payment towards a loan
// @Tags loans
// @Accept json
// @Produce json
// @Param loanId path int true "Loan ID"
// @Param body body service.RepayLoanRequest true "Repayment details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/loans/repay/{loanId} [post]
func (h *LoanHandler) RepayLoan(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	// Get loan ID from path parameter
	loanIDStr := c.Params("loanId")
	loanID, err := strconv.ParseUint(loanIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid loan ID",
		})
	}

	// Parse request body
	var req service.RepayLoanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid request body",
		})
	}

	// Set loan ID from path
	req.LoanID = uint(loanID)

	// Process repayment
	err = h.loanService.RepayLoan(userID, req)
	if err != nil {
		errMsg := err.Error()

		// Handle specific error cases
		status := fiber.StatusBadRequest
		switch errMsg {
		case "loan_not_found":
			status = fiber.StatusNotFound
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "loan_not_found",
			})
		case "not_your_loan":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "not_your_loan",
			})
		case "insufficient_balance":
			return c.Status(status).JSON(fiber.Map{
				"error":   true,
				"message": "insufficient_balance",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "failed to process repayment",
				"details": errMsg,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "repayment processed successfully",
	})
}

// CheckBankruptcy handles GET /api/loans/bankruptcy-check
// @Summary Check bankruptcy status
// @Description Check if user is bankrupt and handle debt collection
// @Tags loans
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/loans/bankruptcy-check [get]
func (h *LoanHandler) CheckBankruptcy(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "unauthorized",
		})
	}

	isBankrupt, err := h.loanService.CheckBankruptcy(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to check bankruptcy",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"is_bankrupt": isBankrupt,
	})
}
