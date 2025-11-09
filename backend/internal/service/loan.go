package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
)

// LoanService provides business logic for loan operations
type LoanService struct {
	db *gorm.DB
}

// NewLoanService creates a new loan service instance
func NewLoanService() *LoanService {
	return &LoanService{
		db: database.GetDB(),
	}
}

// Loan interest rates and limits
const (
	FriendsMaxTotal         = 1000.0 // Max total ever borrowed from friends
	FriendsMaxLoans         = 5      // After 5 loans, friends refuse
	FriendsInterestRate     = 0.0    // 0% interest (friends are generous)
	BankInterestRate        = 0.10   // 10% annual rate
	MicrocreditInterestRate = 2.0    // 200% annual rate (predatory)
)

// TakeLoanRequest represents a loan application request
type TakeLoanRequest struct {
	UserID           uint           `json:"user_id"`
	Amount           float64        `json:"amount"`
	Type             model.LoanType `json:"type"`
	CollateralItemID *uint          `json:"collateral_item_id,omitempty"` // Required for bank loans
}

// TakeLoanResponse represents the response after taking a loan
type TakeLoanResponse struct {
	Loan       model.Loan `json:"loan"`
	NewBalance float64    `json:"new_balance"`
	Message    string     `json:"message"`
}

// RepayLoanRequest represents a loan repayment request
type RepayLoanRequest struct {
	LoanID uint    `json:"loan_id"`
	Amount float64 `json:"amount"`
}

// TakeLoan processes a new loan application
func (s *LoanService) TakeLoan(req TakeLoanRequest) (*TakeLoanResponse, error) {
	// Validate amount
	if req.Amount <= 0 {
		return nil, errors.New("loan amount must be positive")
	}

	// Get user
	var user model.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Calculate total current debt
	summary, err := s.GetLoanSummary(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan summary: %w", err)
	}

	// Type-specific validation
	switch req.Type {
	case model.LoanTypeFriends:
		return s.takeFriendsLoan(req, user, summary)
	case model.LoanTypeBank:
		return s.takeBankLoan(req, user)
	case model.LoanTypeMicrocredit:
		return s.takeMicrocreditLoan(req, user)
	default:
		return nil, errors.New("invalid loan type")
	}
}

// takeFriendsLoan processes a loan from friends
func (s *LoanService) takeFriendsLoan(req TakeLoanRequest, user model.User, summary *model.LoanSummary) (*TakeLoanResponse, error) {
	// Check if user exceeded friend loan limit
	if summary.FriendsLoanCount >= FriendsMaxLoans {
		return nil, errors.New("friends_refused")
	}

	// Check if total borrowed from friends exceeds limit
	if summary.TotalFriendsLoaned+req.Amount > FriendsMaxTotal {
		return nil, errors.New("friends_limit_exceeded")
	}

	// Create loan
	loan := model.Loan{
		UserID:            req.UserID,
		Type:              model.LoanTypeFriends,
		PrincipalAmount:   req.Amount,
		RemainingAmount:   req.Amount, // No interest
		InterestRate:      FriendsInterestRate,
		InterestPerSecond: 0, // Friends don't charge interest
		LastInterestAt:    time.Now(),
	}

	if err := s.db.Create(&loan).Error; err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	// Add money to user balance
	user.Balance += req.Amount
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	return &TakeLoanResponse{
		Loan:       loan,
		NewBalance: user.Balance,
		Message:    "Friends lent you money with no interest!",
	}, nil
}

// takeBankLoan processes a bank loan with collateral
func (s *LoanService) takeBankLoan(req TakeLoanRequest, user model.User) (*TakeLoanResponse, error) {
	// Require collateral
	if req.CollateralItemID == nil {
		return nil, errors.New("collateral_required")
	}

	// Get collateral item
	var userItem model.UserItem
	if err := s.db.Preload("Item").First(&userItem, *req.CollateralItemID).Error; err != nil {
		return nil, errors.New("collateral_not_found")
	}

	// Verify ownership
	if userItem.UserID != req.UserID {
		return nil, errors.New("not_your_item")
	}

	// Verify item is not already collateral
	if userItem.IsCollateral {
		return nil, errors.New("item_already_collateral")
	}

	// Verify collateral value is sufficient (must be >= loan amount)
	// Bank requires car or house as collateral
	if userItem.Item.Type != model.ItemTypeCar && userItem.Item.Type != model.ItemTypeHouse {
		return nil, errors.New("collateral_must_be_car_or_house")
	}

	// Calculate sell value (50% of purchase price)
	collateralValue := userItem.Item.Price * 0.5
	if collateralValue < req.Amount {
		return nil, errors.New("collateral_insufficient")
	}

	// Calculate interest per second from annual rate
	// Formula: principal * (rate / seconds_per_year)
	secondsPerYear := 365.25 * 24 * 60 * 60
	interestPerSecond := req.Amount * (BankInterestRate / secondsPerYear)

	// Create loan
	loan := model.Loan{
		UserID:            req.UserID,
		Type:              model.LoanTypeBank,
		PrincipalAmount:   req.Amount,
		RemainingAmount:   req.Amount,
		InterestRate:      BankInterestRate,
		InterestPerSecond: interestPerSecond,
		CollateralItemID:  req.CollateralItemID,
		LastInterestAt:    time.Now(),
	}

	// Use transaction to ensure atomicity
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create loan
		if err := tx.Create(&loan).Error; err != nil {
			return err
		}

		// Mark item as collateral
		userItem.IsCollateral = true
		if err := tx.Save(&userItem).Error; err != nil {
			return err
		}

		// Update user balance
		user.Balance += req.Amount
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process bank loan: %w", err)
	}

	return &TakeLoanResponse{
		Loan:       loan,
		NewBalance: user.Balance,
		Message:    fmt.Sprintf("Bank loan approved with %s as collateral", userItem.Item.Name),
	}, nil
}

// takeMicrocreditLoan processes a microcredit loan
func (s *LoanService) takeMicrocreditLoan(req TakeLoanRequest, user model.User) (*TakeLoanResponse, error) {
	// Microcredit is available to everyone, no collateral required
	// But charges very high interest

	// Calculate interest per second from annual rate
	secondsPerYear := 365.25 * 24 * 60 * 60
	interestPerSecond := req.Amount * (MicrocreditInterestRate / secondsPerYear)

	// Create loan
	loan := model.Loan{
		UserID:            req.UserID,
		Type:              model.LoanTypeMicrocredit,
		PrincipalAmount:   req.Amount,
		RemainingAmount:   req.Amount,
		InterestRate:      MicrocreditInterestRate,
		InterestPerSecond: interestPerSecond,
		LastInterestAt:    time.Now(),
	}

	// Use transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create loan
		if err := tx.Create(&loan).Error; err != nil {
			return err
		}

		// Update user balance
		user.Balance += req.Amount
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process microcredit: %w", err)
	}

	return &TakeLoanResponse{
		Loan:       loan,
		NewBalance: user.Balance,
		Message:    "Microcredit approved! Be aware of high interest rates.",
	}, nil
}

// GetLoanSummary returns aggregate loan information for a user
func (s *LoanService) GetLoanSummary(userID uint) (*model.LoanSummary, error) {
	// Update all loans interest before calculating
	if err := s.UpdateAllLoansInterest(userID); err != nil {
		return nil, fmt.Errorf("failed to update interest: %w", err)
	}

	var loans []model.Loan
	if err := s.db.Where("user_id = ?", userID).Find(&loans).Error; err != nil {
		return nil, fmt.Errorf("failed to get loans: %w", err)
	}

	summary := &model.LoanSummary{
		TotalDebt:          0,
		InterestPerSecond:  0,
		FriendsLoanCount:   0,
		TotalFriendsLoaned: 0,
		ActiveLoans:        len(loans),
	}

	for _, loan := range loans {
		summary.TotalDebt += loan.RemainingAmount
		summary.InterestPerSecond += loan.InterestPerSecond

		if loan.Type == model.LoanTypeFriends {
			summary.FriendsLoanCount++
			summary.TotalFriendsLoaned += loan.PrincipalAmount
		}
	}

	return summary, nil
}

// GetUserLoans returns all active loans for a user
func (s *LoanService) GetUserLoans(userID uint) ([]model.Loan, error) {
	// Update interest first
	if err := s.UpdateAllLoansInterest(userID); err != nil {
		return nil, fmt.Errorf("failed to update interest: %w", err)
	}

	var loans []model.Loan
	if err := s.db.Preload("CollateralItem.Item").Where("user_id = ?", userID).Find(&loans).Error; err != nil {
		return nil, fmt.Errorf("failed to get loans: %w", err)
	}

	return loans, nil
}

// UpdateAllLoansInterest updates interest for all active loans
func (s *LoanService) UpdateAllLoansInterest(userID uint) error {
	var loans []model.Loan
	if err := s.db.Where("user_id = ?", userID).Find(&loans).Error; err != nil {
		return fmt.Errorf("failed to get loans: %w", err)
	}

	now := time.Now()
	for i := range loans {
		loan := &loans[i]

		// Calculate elapsed time since last interest calculation
		elapsed := now.Sub(loan.LastInterestAt).Seconds()

		// Calculate and add interest
		interestAccrued := loan.InterestPerSecond * elapsed
		loan.RemainingAmount += interestAccrued
		loan.RemainingAmount = math.Round(loan.RemainingAmount*100) / 100 // Round to 2 decimals
		loan.LastInterestAt = now

		// Save updated loan
		if err := s.db.Save(loan).Error; err != nil {
			return fmt.Errorf("failed to update loan interest: %w", err)
		}
	}

	return nil
}

// RepayLoan processes a loan repayment
func (s *LoanService) RepayLoan(userID uint, req RepayLoanRequest) error {
	if req.Amount <= 0 {
		return errors.New("repayment amount must be positive")
	}

	// Get loan
	var loan model.Loan
	if err := s.db.Preload("CollateralItem").First(&loan, req.LoanID).Error; err != nil {
		return errors.New("loan_not_found")
	}

	// Verify ownership
	if loan.UserID != userID {
		return errors.New("not_your_loan")
	}

	// Update interest before repayment
	if err := s.UpdateAllLoansInterest(userID); err != nil {
		return fmt.Errorf("failed to update interest: %w", err)
	}

	// Reload loan to get updated remaining amount
	if err := s.db.Preload("CollateralItem").First(&loan, req.LoanID).Error; err != nil {
		return errors.New("loan_not_found")
	}

	// Get user
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user_not_found")
	}

	// Check if user has enough balance
	if user.Balance < req.Amount {
		return errors.New("insufficient_balance")
	}

	// Cap payment at remaining amount
	paymentAmount := req.Amount
	if paymentAmount > loan.RemainingAmount {
		paymentAmount = loan.RemainingAmount
	}

	// Use transaction
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Deduct from user balance
		user.Balance -= paymentAmount
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		// Reduce loan amount
		loan.RemainingAmount -= paymentAmount
		loan.RemainingAmount = math.Round(loan.RemainingAmount*100) / 100

		// If loan is fully repaid, release collateral if any
		if loan.RemainingAmount <= 0.01 { // Account for floating point precision
			if loan.CollateralItemID != nil && loan.CollateralItem != nil {
				loan.CollateralItem.IsCollateral = false
				if err := tx.Save(loan.CollateralItem).Error; err != nil {
					return err
				}
			}

			// Delete loan record
			if err := tx.Delete(&loan).Error; err != nil {
				return err
			}
		} else {
			// Save updated loan
			if err := tx.Save(&loan).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to process repayment: %w", err)
	}

	return nil
}

// CheckBankruptcy checks if user is bankrupt and handles debt collection
// Returns true if user went bankrupt
func (s *LoanService) CheckBankruptcy(userID uint) (bool, error) {
	// Get user
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	// Get loan summary
	summary, err := s.GetLoanSummary(userID)
	if err != nil {
		return false, fmt.Errorf("failed to get loan summary: %w", err)
	}

	// Check bankruptcy condition: balance <= 0 and has active loans
	if user.Balance <= 0 && summary.ActiveLoans > 0 {
		// User is bankrupt! Collectors come
		if err := s.handleCollectors(userID); err != nil {
			return false, fmt.Errorf("failed to handle collectors: %w", err)
		}
		return true, nil
	}

	return false, nil
}

// handleCollectors handles the debt collection process (removes everything)
func (s *LoanService) handleCollectors(userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete all loans
		if err := tx.Where("user_id = ?", userID).Delete(&model.Loan{}).Error; err != nil {
			return fmt.Errorf("failed to delete loans: %w", err)
		}

		// Release all collateral items (they'll be deleted anyway)
		if err := tx.Model(&model.UserItem{}).Where("user_id = ?", userID).Update("is_collateral", false).Error; err != nil {
			return fmt.Errorf("failed to release collateral: %w", err)
		}

		// Delete all user items
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserItem{}).Error; err != nil {
			return fmt.Errorf("failed to delete items: %w", err)
		}

		// Reset balance to 0
		if err := tx.Model(&model.User{}).Where("id = ?", userID).Update("balance", 0).Error; err != nil {
			return fmt.Errorf("failed to reset balance: %w", err)
		}

		return nil
	})
}
