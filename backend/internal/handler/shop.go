package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smoreg/freezino/backend/internal/service"
)

// ShopHandler handles shop-related HTTP requests
type ShopHandler struct {
	shopService *service.ShopService
}

// NewShopHandler creates a new shop handler instance
func NewShopHandler() *ShopHandler {
	return &ShopHandler{
		shopService: service.NewShopService(),
	}
}

// GetItems handles GET /api/shop/items
// @Summary Get shop items
// @Description Retrieve all shop items with optional filtering by type and rarity
// @Tags shop
// @Accept json
// @Produce json
// @Param type query string false "Item type (clothing, car, house, accessories)"
// @Param rarity query string false "Item rarity (common, rare, epic, legendary)"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/shop/items [get]
func (h *ShopHandler) GetItems(c *fiber.Ctx) error {
	// Get query parameters
	itemType := c.Query("type")
	rarity := c.Query("rarity")

	items, err := h.shopService.GetItems(itemType, rarity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get items",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"count": len(items),
		},
	})
}

// BuyItem handles POST /api/shop/buy/:itemId
// @Summary Buy an item
// @Description Purchase an item from the shop
// @Tags shop
// @Accept json
// @Produce json
// @Param itemId path int true "Item ID"
// @Param user_id query int true "User ID"
// @Success 200 {object} service.BuyItemResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/shop/buy/{itemId} [post]
func (h *ShopHandler) BuyItem(c *fiber.Ctx) error {
	// Get item ID from path parameter
	itemIDStr := c.Params("itemId")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid item ID",
		})
	}

	// Get user ID from query parameter (in production, use JWT)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user_id",
		})
	}

	// Buy the item
	result, err := h.shopService.BuyItem(uint(userID), uint(itemID))
	if err != nil {
		// Check specific error types
		errMsg := err.Error()
		if errMsg == "user not found" || errMsg == "item not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": errMsg,
			})
		}
		if len(errMsg) > 20 && errMsg[:20] == "insufficient balance" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to purchase item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "item purchased successfully",
	})
}

// SellItem handles POST /api/shop/sell/:userItemId
// @Summary Sell an item
// @Description Sell a user's item for 50% of its original price
// @Tags shop
// @Accept json
// @Produce json
// @Param userItemId path int true "User Item ID"
// @Param user_id query int true "User ID"
// @Success 200 {object} service.SellItemResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/shop/sell/{userItemId} [post]
func (h *ShopHandler) SellItem(c *fiber.Ctx) error {
	// Get user item ID from path parameter
	userItemIDStr := c.Params("userItemId")
	userItemID, err := strconv.ParseUint(userItemIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user item ID",
		})
	}

	// Get user ID from query parameter (in production, use JWT)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user_id",
		})
	}

	// Sell the item
	result, err := h.shopService.SellItem(uint(userID), uint(userItemID))
	if err != nil {
		// Check specific error types
		errMsg := err.Error()
		if errMsg == "user not found" || errMsg == "user item not found or does not belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to sell item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "item sold successfully",
	})
}

// GetMyItems handles GET /api/shop/my-items
// @Summary Get user's items
// @Description Retrieve all items owned by the authenticated user
// @Tags shop
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/shop/my-items [get]
func (h *ShopHandler) GetMyItems(c *fiber.Ctx) error {
	// Get user ID from query parameter (in production, use JWT)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user_id",
		})
	}

	items, err := h.shopService.GetMyItems(uint(userID))
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to get user items",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"items": items,
			"count": len(items),
		},
	})
}

// EquipItem handles POST /api/shop/equip/:userItemId
// @Summary Equip an item
// @Description Equip an item (only 1 item per category can be equipped)
// @Tags shop
// @Accept json
// @Produce json
// @Param userItemId path int true "User Item ID"
// @Param user_id query int true "User ID"
// @Success 200 {object} service.UserItemResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/shop/equip/{userItemId} [post]
func (h *ShopHandler) EquipItem(c *fiber.Ctx) error {
	// Get user item ID from path parameter
	userItemIDStr := c.Params("userItemId")
	userItemID, err := strconv.ParseUint(userItemIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user item ID",
		})
	}

	// Get user ID from query parameter (in production, use JWT)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "user_id is required",
		})
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "invalid user_id",
		})
	}

	// Equip the item
	result, err := h.shopService.EquipItem(uint(userID), uint(userItemID))
	if err != nil {
		// Check specific error types
		errMsg := err.Error()
		if errMsg == "user item not found or does not belong to user" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": errMsg,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "failed to equip item",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    result,
		"message": "item equipped successfully",
	})
}
