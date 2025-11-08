package service

import (
	"testing"
	"time"

	"github.com/smoreg/freezino/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTestItem(t *testing.T, db *gorm.DB, name string, itemType model.ItemType, price float64) *model.Item {
	item := &model.Item{
		Name:        name,
		Type:        itemType,
		Price:       price,
		ImageURL:    "https://example.com/item.png",
		Description: "Test item description",
	}

	err := db.Create(item).Error
	require.NoError(t, err, "failed to create test item")
	return item
}

func TestShopServiceGetItems(t *testing.T) {
	db := setupTestDB(t)

	// Create test items
	createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 50.0)
	createTestItem(t, db, "Jeans", model.ItemTypeClothing, 100.0)
	createTestItem(t, db, "Sedan", model.ItemTypeCar, 10000.0)
	createTestItem(t, db, "House", model.ItemTypeHouse, 50000.0)

	service := &ShopService{db: db}

	// Get all items
	items, err := service.GetItems("", "")
	require.NoError(t, err)
	assert.Len(t, items, 4)

	// Items should be sorted by price ascending
	for i := 0; i < len(items)-1; i++ {
		assert.LessOrEqual(t, items[i].Price, items[i+1].Price)
	}
}

func TestShopServiceGetItemsByType(t *testing.T) {
	db := setupTestDB(t)

	// Create test items
	createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 50.0)
	createTestItem(t, db, "Jeans", model.ItemTypeClothing, 100.0)
	createTestItem(t, db, "Sedan", model.ItemTypeCar, 10000.0)

	service := &ShopService{db: db}

	// Get only clothing items
	items, err := service.GetItems(string(model.ItemTypeClothing), "")
	require.NoError(t, err)
	assert.Len(t, items, 2)

	for _, item := range items {
		assert.Equal(t, string(model.ItemTypeClothing), item.Type)
	}
}

func TestShopServiceBuyItem(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)
	item := createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 50.0)

	service := &ShopService{db: db}

	// Buy item
	response, err := service.BuyItem(user.ID, item.ID)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 950.0, response.NewBalance)
	assert.Equal(t, user.ID, response.UserItem.UserID)
	assert.Equal(t, item.ID, response.UserItem.ItemID)
	assert.False(t, response.UserItem.IsEquipped)
	assert.Greater(t, response.TransactionID, uint(0))

	// Verify user balance updated
	var updatedUser model.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 950.0, updatedUser.Balance)

	// Verify user item created
	var userItem model.UserItem
	err = db.Where("user_id = ? AND item_id = ?", user.ID, item.ID).First(&userItem).Error
	require.NoError(t, err)
	assert.Equal(t, user.ID, userItem.UserID)
	assert.Equal(t, item.ID, userItem.ItemID)

	// Verify transaction created
	var transaction model.Transaction
	err = db.Where("user_id = ? AND type = ?", user.ID, model.TransactionTypePurchase).First(&transaction).Error
	require.NoError(t, err)
	assert.Equal(t, -item.Price, transaction.Amount)
	assert.Equal(t, 950.0, transaction.BalanceAfter)
}

func TestShopServiceBuyItemInsufficientBalance(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 10.0) // Not enough money
	item := createTestItem(t, db, "Expensive Item", model.ItemTypeClothing, 1000.0)

	service := &ShopService{db: db}

	// Try to buy item
	_, err := service.BuyItem(user.ID, item.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient balance")

	// Verify balance not changed
	var updatedUser model.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 10.0, updatedUser.Balance)
}

func TestShopServiceBuyItemUserNotFound(t *testing.T) {
	db := setupTestDB(t)
	item := createTestItem(t, db, "Item", model.ItemTypeClothing, 50.0)

	service := &ShopService{db: db}

	_, err := service.BuyItem(9999, item.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestShopServiceBuyItemItemNotFound(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &ShopService{db: db}

	_, err := service.BuyItem(user.ID, 9999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestShopServiceSellItem(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)
	item := createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 100.0)

	// First buy the item
	service := &ShopService{db: db}
	buyResponse, err := service.BuyItem(user.ID, item.ID)
	require.NoError(t, err)

	// Now sell it
	sellResponse, err := service.SellItem(user.ID, buyResponse.UserItem.ID)
	require.NoError(t, err)
	assert.NotNil(t, sellResponse)
	assert.Equal(t, 50.0, sellResponse.SalePrice)   // 50% of 100
	assert.Equal(t, 950.0, sellResponse.NewBalance) // 900 + 50
	assert.Greater(t, sellResponse.TransactionID, uint(0))

	// Verify user balance updated
	var updatedUser model.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 950.0, updatedUser.Balance)

	// Verify user item deleted
	var userItem model.UserItem
	err = db.Where("id = ?", buyResponse.UserItem.ID).First(&userItem).Error
	assert.Error(t, err) // Should not find deleted item

	// Verify transaction created
	var transaction model.Transaction
	err = db.Where("user_id = ? AND type = ?", user.ID, model.TransactionTypeSale).First(&transaction).Error
	require.NoError(t, err)
	assert.Equal(t, 50.0, transaction.Amount)
}

func TestShopServiceSellItemNotOwned(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &ShopService{db: db}

	// Try to sell item that user doesn't own
	_, err := service.SellItem(user.ID, 9999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestShopServiceSellItemWrongUser(t *testing.T) {
	db := setupTestDB(t)
	user1 := createTestUser(t, db, 1000.0)
	user2 := createTestUser(t, db, 1000.0)
	item := createTestItem(t, db, "Item", model.ItemTypeClothing, 100.0)

	service := &ShopService{db: db}

	// User 1 buys item
	buyResponse, err := service.BuyItem(user1.ID, item.ID)
	require.NoError(t, err)

	// User 2 tries to sell user 1's item
	_, err = service.SellItem(user2.ID, buyResponse.UserItem.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestShopServiceGetMyItems(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)
	item1 := createTestItem(t, db, "Item 1", model.ItemTypeClothing, 50.0)
	item2 := createTestItem(t, db, "Item 2", model.ItemTypeCar, 100.0)

	service := &ShopService{db: db}

	// Buy both items
	_, err := service.BuyItem(user.ID, item1.ID)
	require.NoError(t, err)
	time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps
	_, err = service.BuyItem(user.ID, item2.ID)
	require.NoError(t, err)

	// Get my items
	myItems, err := service.GetMyItems(user.ID)
	require.NoError(t, err)
	assert.Len(t, myItems, 2)

	// Items should be sorted by purchased_at DESC (newest first)
	assert.Equal(t, "Item 2", myItems[0].Item.Name)
	assert.Equal(t, "Item 1", myItems[1].Item.Name)
}

func TestShopServiceGetMyItemsEmpty(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &ShopService{db: db}

	myItems, err := service.GetMyItems(user.ID)
	require.NoError(t, err)
	assert.Empty(t, myItems)
}

func TestShopServiceEquipItem(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)
	item := createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 50.0)

	service := &ShopService{db: db}

	// Buy item
	buyResponse, err := service.BuyItem(user.ID, item.ID)
	require.NoError(t, err)
	assert.False(t, buyResponse.UserItem.IsEquipped)

	// Equip item
	equipResponse, err := service.EquipItem(user.ID, buyResponse.UserItem.ID)
	require.NoError(t, err)
	assert.True(t, equipResponse.IsEquipped)

	// Verify in database
	var userItem model.UserItem
	err = db.First(&userItem, buyResponse.UserItem.ID).Error
	require.NoError(t, err)
	assert.True(t, userItem.IsEquipped)
}

func TestShopServiceEquipItemOnlyOnePerCategory(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)
	item1 := createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 50.0)
	item2 := createTestItem(t, db, "Jacket", model.ItemTypeClothing, 100.0)

	service := &ShopService{db: db}

	// Buy both items
	buyResponse1, err := service.BuyItem(user.ID, item1.ID)
	require.NoError(t, err)
	buyResponse2, err := service.BuyItem(user.ID, item2.ID)
	require.NoError(t, err)

	// Equip first item
	_, err = service.EquipItem(user.ID, buyResponse1.UserItem.ID)
	require.NoError(t, err)

	// Equip second item (should unequip first)
	equipResponse2, err := service.EquipItem(user.ID, buyResponse2.UserItem.ID)
	require.NoError(t, err)
	assert.True(t, equipResponse2.IsEquipped)

	// Verify first item is no longer equipped
	var userItem1 model.UserItem
	err = db.First(&userItem1, buyResponse1.UserItem.ID).Error
	require.NoError(t, err)
	assert.False(t, userItem1.IsEquipped, "first item should be unequipped")

	// Verify second item is equipped
	var userItem2 model.UserItem
	err = db.First(&userItem2, buyResponse2.UserItem.ID).Error
	require.NoError(t, err)
	assert.True(t, userItem2.IsEquipped, "second item should be equipped")
}

func TestShopServiceEquipItemDifferentCategories(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 10000.0)
	clothing := createTestItem(t, db, "T-Shirt", model.ItemTypeClothing, 50.0)
	car := createTestItem(t, db, "Car", model.ItemTypeCar, 5000.0)

	service := &ShopService{db: db}

	// Buy both items
	buyResponse1, err := service.BuyItem(user.ID, clothing.ID)
	require.NoError(t, err)
	buyResponse2, err := service.BuyItem(user.ID, car.ID)
	require.NoError(t, err)

	// Equip both items (different categories, both should remain equipped)
	_, err = service.EquipItem(user.ID, buyResponse1.UserItem.ID)
	require.NoError(t, err)
	_, err = service.EquipItem(user.ID, buyResponse2.UserItem.ID)
	require.NoError(t, err)

	// Verify both are equipped
	var userItem1 model.UserItem
	err = db.First(&userItem1, buyResponse1.UserItem.ID).Error
	require.NoError(t, err)
	assert.True(t, userItem1.IsEquipped)

	var userItem2 model.UserItem
	err = db.First(&userItem2, buyResponse2.UserItem.ID).Error
	require.NoError(t, err)
	assert.True(t, userItem2.IsEquipped)
}

func TestShopServiceEquipItemNotOwned(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 1000.0)

	service := &ShopService{db: db}

	_, err := service.EquipItem(user.ID, 9999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestShopServiceTransactionRollback(t *testing.T) {
	db := setupTestDB(t)
	user := createTestUser(t, db, 100.0)
	item := createTestItem(t, db, "Item", model.ItemTypeClothing, 50.0)

	service := &ShopService{db: db}

	// Close the database to simulate failure
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close()

	// Try to buy item (should fail and rollback)
	_, err = service.BuyItem(user.ID, item.ID)
	assert.Error(t, err)

	// Reopen database
	db = setupTestDB(t)
	user = createTestUser(t, db, 100.0)
	item = createTestItem(t, db, "Item", model.ItemTypeClothing, 50.0)

	// Verify balance unchanged (transaction rolled back)
	var checkUser model.User
	err = db.First(&checkUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, 100.0, checkUser.Balance)
}
