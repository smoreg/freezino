package service

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/smoreg/freezino/backend/internal/model"
)

// AuthService handles authentication logic
type AuthService struct {
	db *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required,min=2"`
}

// LoginRequest represents login data
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// giveStarterItems gives new users basic clothing and a cheap house
func (s *AuthService) giveStarterItems(userID uint) error {
	// Find cheapest clothing item
	var clothingItem model.Item
	if err := s.db.Where("type = ?", "clothing").Order("price ASC").First(&clothingItem).Error; err != nil {
		return fmt.Errorf("failed to find starter clothing: %w", err)
	}

	// Find cheapest house item
	var houseItem model.Item
	if err := s.db.Where("type = ?", "house").Order("price ASC").First(&houseItem).Error; err != nil {
		return fmt.Errorf("failed to find starter house: %w", err)
	}

	// Give clothing (equipped)
	clothingUserItem := &model.UserItem{
		UserID:     userID,
		ItemID:     clothingItem.ID,
		IsEquipped: true,
	}
	if err := s.db.Create(clothingUserItem).Error; err != nil {
		return fmt.Errorf("failed to give starter clothing: %w", err)
	}

	// Give house (equipped)
	houseUserItem := &model.UserItem{
		UserID:     userID,
		ItemID:     houseItem.ID,
		IsEquipped: true,
	}
	if err := s.db.Create(houseUserItem).Error; err != nil {
		return fmt.Errorf("failed to give starter house: %w", err)
	}

	return nil
}

// Register creates a new user with username/password
func (s *AuthService) Register(req RegisterRequest) (*model.User, error) {
	// Check if username already exists
	var existing model.User
	if err := s.db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return nil, errors.New("username already taken")
	}

	// Check if email already exists
	if err := s.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		Balance:      1000.00, // Starting balance
		Avatar:       "", // Default empty avatar
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Give starter items
	if err := s.giveStarterItems(user.ID); err != nil {
		// Log error but don't fail registration
		fmt.Printf("Warning: Failed to give starter items to user %d: %v\n", user.ID, err)
	}

	return user, nil
}

// Login authenticates a user with username/password
func (s *AuthService) Login(req LoginRequest) (*model.User, error) {
	var user model.User

	// Find user by username
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user has password (might be Google OAuth user)
	if user.PasswordHash == "" {
		return nil, errors.New("this account uses Google login")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return &user, nil
}
