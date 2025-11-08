package database

import (
	"fmt"
	"log"
	"time"

	"github.com/smoreg/freezino/backend/internal/model"
)

// Seed populates the database with initial data
func Seed() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Seeding database...")

	// Seed test user
	if err := seedTestUser(); err != nil {
		return fmt.Errorf("failed to seed test user: %w", err)
	}

	// Seed items
	if err := seedItems(); err != nil {
		return fmt.Errorf("failed to seed items: %w", err)
	}

	log.Println("Database seeding completed successfully")
	return nil
}

// seedTestUser creates a test user
func seedTestUser() error {
	// Check if test user already exists
	var existingUser model.User
	result := DB.Where("email = ?", "test@freezino.com").First(&existingUser)

	if result.Error == nil {
		log.Println("Test user already exists, skipping...")
		return nil
	}

	testUser := model.User{
		GoogleID:  "test-google-id-123",
		Email:     "test@freezino.com",
		Name:      "Test User",
		Avatar:    "https://via.placeholder.com/150",
		Balance:   1000.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := DB.Create(&testUser).Error; err != nil {
		return err
	}

	log.Printf("Test user created: %s (ID: %d)", testUser.Email, testUser.ID)
	return nil
}

// seedItems creates shop items
func seedItems() error {
	// Check if items already exist
	var count int64
	DB.Model(&model.Item{}).Count(&count)

	if count > 0 {
		log.Printf("Items already exist (%d items), skipping...", count)
		return nil
	}

	items := []model.Item{
		// CLOTHING - Budget (8 items)
		{
			Name:        "Plain T-Shirt",
			Type:        model.ItemTypeClothing,
			Price:       500.00,
			ImageURL:    "https://via.placeholder.com/200x200/808080/FFFFFF?text=Plain+T-Shirt",
			Description: "A simple everyday t-shirt",
		},
		{
			Name:        "Casual Jeans",
			Type:        model.ItemTypeClothing,
			Price:       800.00,
			ImageURL:    "https://via.placeholder.com/200x200/4169E1/FFFFFF?text=Casual+Jeans",
			Description: "Comfortable denim jeans",
		},
		{
			Name:        "Sneakers",
			Type:        model.ItemTypeClothing,
			Price:       1200.00,
			ImageURL:    "https://via.placeholder.com/200x200/FF4500/FFFFFF?text=Sneakers",
			Description: "Everyday sneakers",
		},
		{
			Name:        "Hoodie",
			Type:        model.ItemTypeClothing,
			Price:       1500.00,
			ImageURL:    "https://via.placeholder.com/200x200/696969/FFFFFF?text=Hoodie",
			Description: "Warm and cozy hoodie",
		},

		// CLOTHING - Mid-range (6 items)
		{
			Name:        "Designer Shirt",
			Type:        model.ItemTypeClothing,
			Price:       2000.00,
			ImageURL:    "https://via.placeholder.com/200x200/1E90FF/FFFFFF?text=Designer+Shirt",
			Description: "Stylish designer shirt",
		},
		{
			Name:        "Leather Jacket",
			Type:        model.ItemTypeClothing,
			Price:       3500.00,
			ImageURL:    "https://via.placeholder.com/200x200/000000/FFFFFF?text=Leather+Jacket",
			Description: "Premium leather jacket",
		},
		{
			Name:        "Business Suit",
			Type:        model.ItemTypeClothing,
			Price:       5000.00,
			ImageURL:    "https://via.placeholder.com/200x200/2F4F4F/FFFFFF?text=Business+Suit",
			Description: "Professional business suit",
		},
		{
			Name:        "Designer Dress",
			Type:        model.ItemTypeClothing,
			Price:       4500.00,
			ImageURL:    "https://via.placeholder.com/200x200/FF1493/FFFFFF?text=Designer+Dress",
			Description: "Elegant designer dress",
		},

		// CLOTHING - Luxury (4 items)
		{
			Name:        "Luxury Watch",
			Type:        model.ItemTypeAccessories,
			Price:       15000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FFD700/000000?text=Luxury+Watch",
			Description: "Premium Swiss watch",
		},
		{
			Name:        "Diamond Ring",
			Type:        model.ItemTypeAccessories,
			Price:       25000.00,
			ImageURL:    "https://via.placeholder.com/200x200/B9F2FF/000000?text=Diamond+Ring",
			Description: "Sparkling diamond ring",
		},
		{
			Name:        "Designer Handbag",
			Type:        model.ItemTypeAccessories,
			Price:       8000.00,
			ImageURL:    "https://via.placeholder.com/200x200/8B4513/FFFFFF?text=Designer+Bag",
			Description: "Luxury designer handbag",
		},
		{
			Name:        "Gold Necklace",
			Type:        model.ItemTypeAccessories,
			Price:       12000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FFD700/000000?text=Gold+Necklace",
			Description: "18k gold necklace",
		},

		// CARS - Budget (4 items)
		{
			Name:        "Old Sedan",
			Type:        model.ItemTypeCar,
			Price:       1000.00,
			ImageURL:    "https://via.placeholder.com/300x200/A9A9A9/FFFFFF?text=Old+Sedan",
			Description: "Reliable but worn sedan",
		},
		{
			Name:        "Compact Car",
			Type:        model.ItemTypeCar,
			Price:       5000.00,
			ImageURL:    "https://via.placeholder.com/300x200/4169E1/FFFFFF?text=Compact+Car",
			Description: "Small and fuel-efficient",
		},
		{
			Name:        "Family Sedan",
			Type:        model.ItemTypeCar,
			Price:       12000.00,
			ImageURL:    "https://via.placeholder.com/300x200/2F4F4F/FFFFFF?text=Family+Sedan",
			Description: "Spacious family car",
		},
		{
			Name:        "Used SUV",
			Type:        model.ItemTypeCar,
			Price:       18000.00,
			ImageURL:    "https://via.placeholder.com/300x200/228B22/FFFFFF?text=Used+SUV",
			Description: "Pre-owned SUV in good condition",
		},

		// CARS - Mid-range (4 items)
		{
			Name:        "New SUV",
			Type:        model.ItemTypeCar,
			Price:       35000.00,
			ImageURL:    "https://via.placeholder.com/300x200/006400/FFFFFF?text=New+SUV",
			Description: "Brand new SUV",
		},
		{
			Name:        "Sports Coupe",
			Type:        model.ItemTypeCar,
			Price:       45000.00,
			ImageURL:    "https://via.placeholder.com/300x200/DC143C/FFFFFF?text=Sports+Coupe",
			Description: "Fast and stylish coupe",
		},
		{
			Name:        "Luxury Sedan",
			Type:        model.ItemTypeCar,
			Price:       60000.00,
			ImageURL:    "https://via.placeholder.com/300x200/191970/FFFFFF?text=Luxury+Sedan",
			Description: "High-end luxury sedan",
		},
		{
			Name:        "Electric Car",
			Type:        model.ItemTypeCar,
			Price:       55000.00,
			ImageURL:    "https://via.placeholder.com/300x200/00CED1/FFFFFF?text=Electric+Car",
			Description: "Modern electric vehicle",
		},

		// CARS - Luxury (4 items)
		{
			Name:        "Tesla Model S",
			Type:        model.ItemTypeCar,
			Price:       95000.00,
			ImageURL:    "https://via.placeholder.com/300x200/1E90FF/FFFFFF?text=Tesla+Model+S",
			Description: "Premium electric luxury sedan",
		},
		{
			Name:        "Porsche 911",
			Type:        model.ItemTypeCar,
			Price:       125000.00,
			ImageURL:    "https://via.placeholder.com/300x200/FFD700/000000?text=Porsche+911",
			Description: "Iconic sports car",
		},
		{
			Name:        "Mercedes S-Class",
			Type:        model.ItemTypeCar,
			Price:       110000.00,
			ImageURL:    "https://via.placeholder.com/300x200/000000/FFFFFF?text=Mercedes+S",
			Description: "Ultimate luxury sedan",
		},
		{
			Name:        "Ferrari F8",
			Type:        model.ItemTypeCar,
			Price:       280000.00,
			ImageURL:    "https://via.placeholder.com/300x200/FF0000/FFFFFF?text=Ferrari+F8",
			Description: "Italian supercar",
		},

		// HOUSES - Budget (3 items)
		{
			Name:        "Studio Apartment",
			Type:        model.ItemTypeHouse,
			Price:       5000.00,
			ImageURL:    "https://via.placeholder.com/400x300/778899/FFFFFF?text=Studio",
			Description: "Cozy studio apartment",
		},
		{
			Name:        "Small Apartment",
			Type:        model.ItemTypeHouse,
			Price:       15000.00,
			ImageURL:    "https://via.placeholder.com/400x300/4682B4/FFFFFF?text=Small+Apt",
			Description: "One-bedroom apartment",
		},
		{
			Name:        "Suburban House",
			Type:        model.ItemTypeHouse,
			Price:       35000.00,
			ImageURL:    "https://via.placeholder.com/400x300/8FBC8F/FFFFFF?text=Suburban",
			Description: "Nice house in the suburbs",
		},

		// HOUSES - Mid-range (3 items)
		{
			Name:        "City Condo",
			Type:        model.ItemTypeHouse,
			Price:       75000.00,
			ImageURL:    "https://via.placeholder.com/400x300/708090/FFFFFF?text=City+Condo",
			Description: "Modern downtown condo",
		},
		{
			Name:        "Family Home",
			Type:        model.ItemTypeHouse,
			Price:       120000.00,
			ImageURL:    "https://via.placeholder.com/400x300/CD853F/FFFFFF?text=Family+Home",
			Description: "Spacious family home",
		},
		{
			Name:        "Beach House",
			Type:        model.ItemTypeHouse,
			Price:       200000.00,
			ImageURL:    "https://via.placeholder.com/400x300/87CEEB/FFFFFF?text=Beach+House",
			Description: "Beautiful beachfront property",
		},

		// HOUSES - Luxury (4 items)
		{
			Name:        "Luxury Penthouse",
			Type:        model.ItemTypeHouse,
			Price:       500000.00,
			ImageURL:    "https://via.placeholder.com/400x300/2F4F4F/FFFFFF?text=Penthouse",
			Description: "Top-floor luxury penthouse",
		},
		{
			Name:        "Modern Mansion",
			Type:        model.ItemTypeHouse,
			Price:       750000.00,
			ImageURL:    "https://via.placeholder.com/400x300/DAA520/FFFFFF?text=Mansion",
			Description: "Stunning modern mansion",
		},
		{
			Name:        "Private Estate",
			Type:        model.ItemTypeHouse,
			Price:       1000000.00,
			ImageURL:    "https://via.placeholder.com/400x300/8B4513/FFFFFF?text=Estate",
			Description: "Exclusive private estate",
		},
		{
			Name:        "Island Villa",
			Type:        model.ItemTypeHouse,
			Price:       1000000.00,
			ImageURL:    "https://via.placeholder.com/400x300/20B2AA/FFFFFF?text=Island+Villa",
			Description: "Private island villa paradise",
		},
	}

	// Create all items
	for _, item := range items {
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		if err := DB.Create(&item).Error; err != nil {
			return err
		}
	}

	log.Printf("Created %d shop items", len(items))
	return nil
}

// ClearData removes all data from the database (keeps schema)
func ClearData() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Clearing all data...")

	// Delete in reverse order of dependencies
	if err := DB.Exec("DELETE FROM game_sessions").Error; err != nil {
		return err
	}
	if err := DB.Exec("DELETE FROM work_sessions").Error; err != nil {
		return err
	}
	if err := DB.Exec("DELETE FROM user_items").Error; err != nil {
		return err
	}
	if err := DB.Exec("DELETE FROM transactions").Error; err != nil {
		return err
	}
	if err := DB.Exec("DELETE FROM items").Error; err != nil {
		return err
	}
	if err := DB.Exec("DELETE FROM users").Error; err != nil {
		return err
	}

	log.Println("All data cleared successfully")
	return nil
}
