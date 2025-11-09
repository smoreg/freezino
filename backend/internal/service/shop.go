package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/smoreg/freezino/backend/internal/database"
	"github.com/smoreg/freezino/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ShopService provides business logic for shop operations
type ShopService struct {
	db *gorm.DB
}

// NewShopService creates a new shop service instance
func NewShopService() *ShopService {
	return &ShopService{
		db: database.GetDB(),
	}
}

// ItemResponse represents a shop item response
type ItemResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// UserItemResponse represents a user's item with full details
type UserItemResponse struct {
	ID          uint         `json:"id"`
	UserID      uint         `json:"user_id"`
	ItemID      uint         `json:"item_id"`
	PurchasedAt string       `json:"purchased_at"`
	IsEquipped  bool         `json:"is_equipped"`
	Item        ItemResponse `json:"item"`
}

// BuyItemResponse represents the response after buying an item
type BuyItemResponse struct {
	UserItem      UserItemResponse `json:"user_item"`
	NewBalance    float64          `json:"new_balance"`
	TransactionID uint             `json:"transaction_id"`
}

// SellItemResponse represents the response after selling an item
type SellItemResponse struct {
	SalePrice     float64 `json:"sale_price"`
	NewBalance    float64 `json:"new_balance"`
	TransactionID uint    `json:"transaction_id"`
}

// GetItems retrieves shop items with optional filtering
func (s *ShopService) GetItems(itemType string, rarity string) ([]ItemResponse, error) {
	var items []model.Item
	query := s.db.Model(&model.Item{})

	// Filter by type if provided
	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}

	// Filter by rarity if provided (future feature - Rarity field doesn't exist yet)
	// This will be implemented when Клод 20 adds the rarity field
	if rarity != "" {
		// query = query.Where("rarity = ?", rarity)
		// For now, we'll ignore rarity filter as the field doesn't exist yet
	}

	// Order by price ascending
	query = query.Order("price ASC")

	if err := query.Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	// Convert to response format
	response := make([]ItemResponse, len(items))
	for i, item := range items {
		response[i] = ItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Type:        string(item.Type),
			Price:       item.Price,
			ImageURL:    item.ImageURL,
			Description: item.Description,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response, nil
}

// BuyItem handles purchasing an item
func (s *ShopService) BuyItem(userID uint, itemID uint) (*BuyItemResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user with lock to prevent race conditions
	var user model.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get item
	var item model.Item
	if err := tx.First(&item, itemID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	// Check if user has enough balance
	if user.Balance < item.Price {
		tx.Rollback()
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", user.Balance, item.Price)
	}

	// Deduct price from user balance
	newBalance := user.Balance - item.Price
	if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Create user item
	userItem := model.UserItem{
		UserID:      userID,
		ItemID:      itemID,
		PurchasedAt: time.Now(),
		IsEquipped:  false,
	}
	if err := tx.Create(&userItem).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create user item: %w", err)
	}

	// Create transaction record
	transaction := model.Transaction{
		UserID:       userID,
		Type:         model.TransactionTypePurchase,
		Amount:       -item.Price, // Negative because it's a purchase
		BalanceAfter: newBalance,
		Description:  fmt.Sprintf("Purchased %s", item.Name),
		CreatedAt:    time.Now(),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Build response
	response := &BuyItemResponse{
		UserItem: UserItemResponse{
			ID:          userItem.ID,
			UserID:      userItem.UserID,
			ItemID:      userItem.ItemID,
			PurchasedAt: userItem.PurchasedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsEquipped:  userItem.IsEquipped,
			Item: ItemResponse{
				ID:          item.ID,
				Name:        item.Name,
				Type:        string(item.Type),
				Price:       item.Price,
				ImageURL:    item.ImageURL,
				Description: item.Description,
				CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		},
		NewBalance:    newBalance,
		TransactionID: transaction.ID,
	}

	return response, nil
}

// SellItem handles selling a user's item
func (s *ShopService) SellItem(userID uint, userItemID uint) (*SellItemResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user with lock
	var user model.User
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get user item with item details
	var userItem model.UserItem
	if err := tx.Preload("Item").Where("id = ? AND user_id = ?", userItemID, userID).First(&userItem).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user item not found or does not belong to user")
		}
		return nil, fmt.Errorf("failed to get user item: %w", err)
	}

	// Calculate sale price (50% of original price)
	salePrice := userItem.Item.Price * 0.5
	newBalance := user.Balance + salePrice

	// Update user balance
	if err := tx.Model(&user).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Delete user item
	if err := tx.Delete(&userItem).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete user item: %w", err)
	}

	// Create transaction record
	transaction := model.Transaction{
		UserID:       userID,
		Type:         model.TransactionTypeSale,
		Amount:       salePrice, // Positive because it's income
		BalanceAfter: newBalance,
		Description:  fmt.Sprintf("Sold %s", userItem.Item.Name),
		CreatedAt:    time.Now(),
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	response := &SellItemResponse{
		SalePrice:     salePrice,
		NewBalance:    newBalance,
		TransactionID: transaction.ID,
	}

	return response, nil
}

// GetMyItems retrieves all items owned by a user
func (s *ShopService) GetMyItems(userID uint) ([]UserItemResponse, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Get user items with item details
	var userItems []model.UserItem
	if err := s.db.Preload("Item").Where("user_id = ?", userID).Order("purchased_at DESC").Find(&userItems).Error; err != nil {
		return nil, fmt.Errorf("failed to get user items: %w", err)
	}

	// Convert to response format
	response := make([]UserItemResponse, len(userItems))
	for i, userItem := range userItems {
		response[i] = UserItemResponse{
			ID:          userItem.ID,
			UserID:      userItem.UserID,
			ItemID:      userItem.ItemID,
			PurchasedAt: userItem.PurchasedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsEquipped:  userItem.IsEquipped,
			Item: ItemResponse{
				ID:          userItem.Item.ID,
				Name:        userItem.Item.Name,
				Type:        string(userItem.Item.Type),
				Price:       userItem.Item.Price,
				ImageURL:    userItem.Item.ImageURL,
				Description: userItem.Item.Description,
				CreatedAt:   userItem.Item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			},
		}
	}

	return response, nil
}

// EquipItem equips an item (only 1 item per category can be equipped)
func (s *ShopService) EquipItem(userID uint, userItemID uint) (*UserItemResponse, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user item with item details
	var userItem model.UserItem
	if err := tx.Preload("Item").Where("id = ? AND user_id = ?", userItemID, userID).First(&userItem).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user item not found or does not belong to user")
		}
		return nil, fmt.Errorf("failed to get user item: %w", err)
	}

	// Unequip all items of the same type for this user
	// First, get all user items of the same type
	var sameTypeUserItems []model.UserItem
	if err := tx.Joins("Item").
		Where("user_items.user_id = ? AND user_items.id != ? AND Item.type = ?", userID, userItemID, userItem.Item.Type).
		Find(&sameTypeUserItems).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find other items: %w", err)
	}

	// Unequip them
	for _, item := range sameTypeUserItems {
		if err := tx.Model(&item).Update("is_equipped", false).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to unequip other items: %w", err)
		}
	}

	// Equip the selected item
	if err := tx.Model(&userItem).Update("is_equipped", true).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to equip item: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Reload to get updated state
	if err := s.db.Preload("Item").First(&userItem, userItemID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload user item: %w", err)
	}

	response := &UserItemResponse{
		ID:          userItem.ID,
		UserID:      userItem.UserID,
		ItemID:      userItem.ItemID,
		PurchasedAt: userItem.PurchasedAt.Format("2006-01-02T15:04:05Z07:00"),
		IsEquipped:  userItem.IsEquipped,
		Item: ItemResponse{
			ID:          userItem.Item.ID,
			Name:        userItem.Item.Name,
			Type:        string(userItem.Item.Type),
			Price:       userItem.Item.Price,
			ImageURL:    userItem.Item.ImageURL,
			Description: userItem.Item.Description,
			CreatedAt:   userItem.Item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}

	return response, nil
}
