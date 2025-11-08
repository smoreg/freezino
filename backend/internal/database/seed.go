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
		// CLOTHING - Common (5 items, $500-$2,000)
		{
			Name:        "Plain T-Shirt",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       500.00,
			ImageURL:    "https://via.placeholder.com/200x200/808080/FFFFFF?text=Plain+T-Shirt",
			Description: "A simple everyday t-shirt",
		},
		{
			Name:        "Casual Jeans",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       800.00,
			ImageURL:    "https://via.placeholder.com/200x200/4169E1/FFFFFF?text=Casual+Jeans",
			Description: "Comfortable denim jeans",
		},
		{
			Name:        "Sneakers",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       1200.00,
			ImageURL:    "https://via.placeholder.com/200x200/FF4500/FFFFFF?text=Sneakers",
			Description: "Everyday sneakers",
		},
		{
			Name:        "Hoodie",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       1500.00,
			ImageURL:    "https://via.placeholder.com/200x200/696969/FFFFFF?text=Hoodie",
			Description: "Warm and cozy hoodie",
		},
		{
			Name:        "Designer Shirt",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       2000.00,
			ImageURL:    "https://via.placeholder.com/200x200/1E90FF/FFFFFF?text=Designer+Shirt",
			Description: "Stylish designer shirt",
		},

		// CLOTHING - Rare (6 items, $3,500-$8,000)
		{
			Name:        "Leather Jacket",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       3500.00,
			ImageURL:    "https://via.placeholder.com/200x200/000000/FFFFFF?text=Leather+Jacket",
			Description: "Premium leather jacket",
		},
		{
			Name:        "Designer Dress",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       4500.00,
			ImageURL:    "https://via.placeholder.com/200x200/FF1493/FFFFFF?text=Designer+Dress",
			Description: "Elegant designer dress",
		},
		{
			Name:        "Business Suit",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       5000.00,
			ImageURL:    "https://via.placeholder.com/200x200/2F4F4F/FFFFFF?text=Business+Suit",
			Description: "Professional business suit",
		},
		{
			Name:        "Evening Gown",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       6500.00,
			ImageURL:    "https://via.placeholder.com/200x200/9400D3/FFFFFF?text=Evening+Gown",
			Description: "Glamorous evening gown",
		},
		{
			Name:        "Tuxedo",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       7000.00,
			ImageURL:    "https://via.placeholder.com/200x200/000000/FFFFFF?text=Tuxedo",
			Description: "Classic black tuxedo",
		},
		{
			Name:        "Designer Coat",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       8000.00,
			ImageURL:    "https://via.placeholder.com/200x200/8B4513/FFFFFF?text=Designer+Coat",
			Description: "High-fashion winter coat",
		},

		// CLOTHING - Epic (3 items, $15,000-$35,000)
		{
			Name:        "Custom Tailored Suit",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       15000.00,
			ImageURL:    "https://via.placeholder.com/200x200/4B0082/FFFFFF?text=Custom+Suit",
			Description: "Hand-tailored bespoke suit",
		},
		{
			Name:        "Haute Couture Dress",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       25000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FF69B4/FFFFFF?text=Haute+Couture",
			Description: "Exclusive haute couture piece",
		},
		{
			Name:        "Luxury Fur Coat",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       35000.00,
			ImageURL:    "https://via.placeholder.com/200x200/D2691E/FFFFFF?text=Fur+Coat",
			Description: "Premium fur coat",
		},

		// CLOTHING - Legendary (1 item, $50,000)
		{
			Name:        "Limited Edition Designer Collection",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityLegendary,
			Price:       50000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FFD700/000000?text=Limited+Edition",
			Description: "Rare runway piece from exclusive collection",
		},

		// ACCESSORIES - Common (5 items, $500-$3,000)
		{
			Name:        "Sunglasses",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       500.00,
			ImageURL:    "https://via.placeholder.com/200x200/000000/FFFFFF?text=Sunglasses",
			Description: "Stylish sunglasses",
		},
		{
			Name:        "Leather Wallet",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       800.00,
			ImageURL:    "https://via.placeholder.com/200x200/8B4513/FFFFFF?text=Wallet",
			Description: "Quality leather wallet",
		},
		{
			Name:        "Casual Watch",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       1500.00,
			ImageURL:    "https://via.placeholder.com/200x200/C0C0C0/000000?text=Casual+Watch",
			Description: "Everyday wristwatch",
		},
		{
			Name:        "Belt",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       600.00,
			ImageURL:    "https://via.placeholder.com/200x200/654321/FFFFFF?text=Belt",
			Description: "Classic leather belt",
		},
		{
			Name:        "Backpack",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       1200.00,
			ImageURL:    "https://via.placeholder.com/200x200/2F4F4F/FFFFFF?text=Backpack",
			Description: "Practical everyday backpack",
		},

		// ACCESSORIES - Rare (6 items, $5,000-$15,000)
		{
			Name:        "Designer Handbag",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       8000.00,
			ImageURL:    "https://via.placeholder.com/200x200/8B4513/FFFFFF?text=Designer+Bag",
			Description: "Luxury designer handbag",
		},
		{
			Name:        "Gold Necklace",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       12000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FFD700/000000?text=Gold+Necklace",
			Description: "18k gold necklace",
		},
		{
			Name:        "Designer Briefcase",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       5500.00,
			ImageURL:    "https://via.placeholder.com/200x200/2F4F4F/FFFFFF?text=Briefcase",
			Description: "Premium leather briefcase",
		},
		{
			Name:        "Pearl Earrings",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       9000.00,
			ImageURL:    "https://via.placeholder.com/200x200/F5F5DC/000000?text=Pearl+Earrings",
			Description: "Natural pearl drop earrings",
		},
		{
			Name:        "Silver Bracelet",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       6500.00,
			ImageURL:    "https://via.placeholder.com/200x200/C0C0C0/000000?text=Bracelet",
			Description: "Sterling silver bracelet",
		},
		{
			Name:        "Luxury Watch",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       15000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FFD700/000000?text=Luxury+Watch",
			Description: "Premium Swiss watch",
		},

		// ACCESSORIES - Epic (2 items, $20,000-$25,000)
		{
			Name:        "Diamond Ring",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityEpic,
			Price:       25000.00,
			ImageURL:    "https://via.placeholder.com/200x200/B9F2FF/000000?text=Diamond+Ring",
			Description: "Sparkling diamond ring",
		},
		{
			Name:        "Platinum Cufflinks",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityEpic,
			Price:       20000.00,
			ImageURL:    "https://via.placeholder.com/200x200/E5E4E2/000000?text=Cufflinks",
			Description: "Handcrafted platinum cufflinks",
		},

		// ACCESSORIES - Legendary (2 items, $50,000+)
		{
			Name:        "Rare Collectible Watch",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityLegendary,
			Price:       85000.00,
			ImageURL:    "https://via.placeholder.com/200x200/FFD700/000000?text=Collectible+Watch",
			Description: "Limited edition timepiece from prestigious watchmaker",
		},
		{
			Name:        "Diamond Necklace Set",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityLegendary,
			Price:       120000.00,
			ImageURL:    "https://via.placeholder.com/200x200/B9F2FF/000000?text=Diamond+Set",
			Description: "Exquisite diamond necklace and earring set",
		},

		// CARS - Common (2 items, $1,000-$5,000)
		{
			Name:        "Old Sedan",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       1000.00,
			ImageURL:    "https://via.placeholder.com/300x200/A9A9A9/FFFFFF?text=Old+Sedan",
			Description: "Reliable but worn sedan",
		},
		{
			Name:        "Compact Car",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       5000.00,
			ImageURL:    "https://via.placeholder.com/300x200/4169E1/FFFFFF?text=Compact+Car",
			Description: "Small and fuel-efficient",
		},

		// CARS - Rare (4 items, $12,000-$45,000)
		{
			Name:        "Family Sedan",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       12000.00,
			ImageURL:    "https://via.placeholder.com/300x200/2F4F4F/FFFFFF?text=Family+Sedan",
			Description: "Spacious family car",
		},
		{
			Name:        "Used SUV",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       18000.00,
			ImageURL:    "https://via.placeholder.com/300x200/228B22/FFFFFF?text=Used+SUV",
			Description: "Pre-owned SUV in good condition",
		},
		{
			Name:        "New SUV",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       35000.00,
			ImageURL:    "https://via.placeholder.com/300x200/006400/FFFFFF?text=New+SUV",
			Description: "Brand new SUV",
		},
		{
			Name:        "Sports Coupe",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       45000.00,
			ImageURL:    "https://via.placeholder.com/300x200/DC143C/FFFFFF?text=Sports+Coupe",
			Description: "Fast and stylish coupe",
		},

		// CARS - Epic (3 items, $55,000-$95,000)
		{
			Name:        "Electric Car",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       55000.00,
			ImageURL:    "https://via.placeholder.com/300x200/00CED1/FFFFFF?text=Electric+Car",
			Description: "Modern electric vehicle",
		},
		{
			Name:        "Luxury Sedan",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       60000.00,
			ImageURL:    "https://via.placeholder.com/300x200/191970/FFFFFF?text=Luxury+Sedan",
			Description: "High-end luxury sedan",
		},
		{
			Name:        "Tesla Model S",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       95000.00,
			ImageURL:    "https://via.placeholder.com/300x200/1E90FF/FFFFFF?text=Tesla+Model+S",
			Description: "Premium electric luxury sedan",
		},

		// CARS - Legendary (3 items, $110,000-$500,000)
		{
			Name:        "Mercedes S-Class",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       110000.00,
			ImageURL:    "https://via.placeholder.com/300x200/000000/FFFFFF?text=Mercedes+S",
			Description: "Ultimate luxury sedan",
		},
		{
			Name:        "Porsche 911",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       125000.00,
			ImageURL:    "https://via.placeholder.com/300x200/FFD700/000000?text=Porsche+911",
			Description: "Iconic sports car",
		},
		{
			Name:        "Ferrari F8",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       280000.00,
			ImageURL:    "https://via.placeholder.com/300x200/FF0000/FFFFFF?text=Ferrari+F8",
			Description: "Italian supercar",
		},
		{
			Name:        "Lamborghini Aventador",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       500000.00,
			ImageURL:    "https://via.placeholder.com/300x200/FFA500/000000?text=Lamborghini",
			Description: "Legendary Italian supercar",
		},

		// HOUSES - Common (2 items, $2,000-$5,000)
		{
			Name:        "Studio Apartment",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       5000.00,
			ImageURL:    "https://via.placeholder.com/400x300/778899/FFFFFF?text=Studio",
			Description: "Cozy studio apartment",
		},
		{
			Name:        "Small Apartment",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       15000.00,
			ImageURL:    "https://via.placeholder.com/400x300/4682B4/FFFFFF?text=Small+Apt",
			Description: "One-bedroom apartment",
		},

		// HOUSES - Rare (4 items, $35,000-$120,000)
		{
			Name:        "Suburban House",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       35000.00,
			ImageURL:    "https://via.placeholder.com/400x300/8FBC8F/FFFFFF?text=Suburban",
			Description: "Nice house in the suburbs",
		},
		{
			Name:        "City Condo",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       75000.00,
			ImageURL:    "https://via.placeholder.com/400x300/708090/FFFFFF?text=City+Condo",
			Description: "Modern downtown condo",
		},
		{
			Name:        "Family Home",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       120000.00,
			ImageURL:    "https://via.placeholder.com/400x300/CD853F/FFFFFF?text=Family+Home",
			Description: "Spacious family home",
		},
		{
			Name:        "Lake House",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       95000.00,
			ImageURL:    "https://via.placeholder.com/400x300/4682B4/FFFFFF?text=Lake+House",
			Description: "Peaceful lakeside retreat",
		},

		// HOUSES - Epic (2 items, $200,000-$500,000)
		{
			Name:        "Beach House",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityEpic,
			Price:       200000.00,
			ImageURL:    "https://via.placeholder.com/400x300/87CEEB/FFFFFF?text=Beach+House",
			Description: "Beautiful beachfront property",
		},
		{
			Name:        "Luxury Penthouse",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityEpic,
			Price:       500000.00,
			ImageURL:    "https://via.placeholder.com/400x300/2F4F4F/FFFFFF?text=Penthouse",
			Description: "Top-floor luxury penthouse",
		},

		// HOUSES - Legendary (3 items, $750,000-$1,000,000)
		{
			Name:        "Modern Mansion",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
			Price:       750000.00,
			ImageURL:    "https://via.placeholder.com/400x300/DAA520/FFFFFF?text=Mansion",
			Description: "Stunning modern mansion",
		},
		{
			Name:        "Private Estate",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
			Price:       1000000.00,
			ImageURL:    "https://via.placeholder.com/400x300/8B4513/FFFFFF?text=Estate",
			Description: "Exclusive private estate",
		},
		{
			Name:        "Island Villa",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
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
