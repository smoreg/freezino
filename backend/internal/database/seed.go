package database

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

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

// seedTestUser creates test users with username/password
func seedTestUser() error {
	// Define test users
	testUsers := []struct {
		username string
		email    string
		name     string
		password string
	}{
		{
			username: "aaa",
			email:    "aaa@test.com",
			name:     "Test User AAA",
			password: "aaa",
		},
		{
			username: "testuser123",
			email:    "testuser123@test.com",
			name:     "Test User 123",
			password: "testuser123",
		},
	}

	for _, userData := range testUsers {
		// Check if user already exists
		var existingUser model.User
		result := DB.Where("username = ?", userData.username).First(&existingUser)

		if result.Error == nil {
			log.Printf("Test user '%s' already exists, skipping...", userData.username)
			continue
		}

		// Hash password
		hashedPassword, err := hashPassword(userData.password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", userData.username, err)
		}

		testUser := model.User{
			Username:     userData.username,
			Email:        userData.email,
			Name:         userData.name,
			PasswordHash: hashedPassword,
			Avatar:       fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", userData.username),
			Balance:      1000.00,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := DB.Create(&testUser).Error; err != nil {
			return fmt.Errorf("failed to create test user %s: %w", userData.username, err)
		}

		log.Printf("Test user created: %s (ID: %d, username: %s)", testUser.Email, testUser.ID, testUser.Username)
	}

	return nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
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
			ImageURL:    "/images/clothing/plain-tshirt.jpg",
			Description: "A simple everyday t-shirt",
		},
		{
			Name:        "Casual Jeans",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       800.00,
			ImageURL:    "/images/clothing/casual-jeans.jpg",
			Description: "Comfortable denim jeans",
		},
		{
			Name:        "Sneakers",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       1200.00,
			ImageURL:    "/images/clothing/sneakers.jpg",
			Description: "Everyday sneakers",
		},
		{
			Name:        "Hoodie",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       1500.00,
			ImageURL:    "/images/clothing/hoodie.jpg",
			Description: "Warm and cozy hoodie",
		},
		{
			Name:        "Designer Shirt",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       2000.00,
			ImageURL:    "/images/clothing/designer-shirt.jpg",
			Description: "Stylish designer shirt",
		},

		// CLOTHING - Rare (6 items, $3,500-$8,000)
		{
			Name:        "Leather Jacket",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       3500.00,
			ImageURL:    "/images/clothing/leather-jacket.jpg",
			Description: "Premium leather jacket",
		},
		{
			Name:        "Designer Dress",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       4500.00,
			ImageURL:    "/images/clothing/designer-dress.jpg",
			Description: "Elegant designer dress",
		},
		{
			Name:        "Business Suit",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       5000.00,
			ImageURL:    "/images/clothing/business-suit.jpg",
			Description: "Professional business suit",
		},
		{
			Name:        "Evening Gown",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       6500.00,
			ImageURL:    "/images/clothing/evening-gown.jpg",
			Description: "Glamorous evening gown",
		},
		{
			Name:        "Tuxedo",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       7000.00,
			ImageURL:    "/images/clothing/tuxedo.jpg",
			Description: "Classic black tuxedo",
		},
		{
			Name:        "Designer Coat",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       8000.00,
			ImageURL:    "/images/clothing/designer-coat.jpg",
			Description: "High-fashion winter coat",
		},

		// CLOTHING - Epic (3 items, $15,000-$35,000)
		{
			Name:        "Custom Tailored Suit",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       15000.00,
			ImageURL:    "/images/clothing/custom-suit.jpg",
			Description: "Hand-tailored bespoke suit",
		},
		{
			Name:        "Haute Couture Dress",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       25000.00,
			ImageURL:    "/images/clothing/haute-couture.jpg",
			Description: "Exclusive haute couture piece",
		},
		{
			Name:        "Luxury Fur Coat",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       35000.00,
			ImageURL:    "/images/clothing/fur-coat.jpg",
			Description: "Premium fur coat",
		},

		// CLOTHING - Legendary (1 item, $50,000)
		{
			Name:        "Limited Edition Designer Collection",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityLegendary,
			Price:       50000.00,
			ImageURL:    "/images/clothing/limited-edition.jpg",
			Description: "Rare runway piece from exclusive collection",
		},

		// ACCESSORIES - Common (5 items, $500-$3,000)
		{
			Name:        "Sunglasses",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       500.00,
			ImageURL:    "/images/accessories/sunglasses.jpg",
			Description: "Stylish sunglasses",
		},
		{
			Name:        "Leather Wallet",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       800.00,
			ImageURL:    "/images/accessories/wallet.jpg",
			Description: "Quality leather wallet",
		},
		{
			Name:        "Casual Watch",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       1500.00,
			ImageURL:    "/images/accessories/casual-watch.jpg",
			Description: "Everyday wristwatch",
		},
		{
			Name:        "Belt",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       600.00,
			ImageURL:    "/images/accessories/belt.jpg",
			Description: "Classic leather belt",
		},
		{
			Name:        "Backpack",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       1200.00,
			ImageURL:    "/images/accessories/backpack.jpg",
			Description: "Practical everyday backpack",
		},

		// ACCESSORIES - Rare (6 items, $5,000-$15,000)
		{
			Name:        "Designer Handbag",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       8000.00,
			ImageURL:    "/images/accessories/designer-handbag.jpg",
			Description: "Luxury designer handbag",
		},
		{
			Name:        "Gold Necklace",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       12000.00,
			ImageURL:    "/images/accessories/gold-necklace.jpg",
			Description: "18k gold necklace",
		},
		{
			Name:        "Designer Briefcase",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       5500.00,
			ImageURL:    "/images/accessories/briefcase.jpg",
			Description: "Premium leather briefcase",
		},
		{
			Name:        "Pearl Earrings",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       9000.00,
			ImageURL:    "/images/accessories/pearl-earrings.jpg",
			Description: "Natural pearl drop earrings",
		},
		{
			Name:        "Silver Bracelet",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       6500.00,
			ImageURL:    "/images/accessories/silver-bracelet.jpg",
			Description: "Sterling silver bracelet",
		},
		{
			Name:        "Luxury Watch",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       15000.00,
			ImageURL:    "/images/accessories/luxury-watch.jpg",
			Description: "Premium Swiss watch",
		},

		// ACCESSORIES - Epic (2 items, $20,000-$25,000)
		{
			Name:        "Diamond Ring",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityEpic,
			Price:       25000.00,
			ImageURL:    "/images/accessories/diamond-ring.jpg",
			Description: "Sparkling diamond ring",
		},
		{
			Name:        "Platinum Cufflinks",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityEpic,
			Price:       20000.00,
			ImageURL:    "/images/accessories/cufflinks.jpg",
			Description: "Handcrafted platinum cufflinks",
		},

		// ACCESSORIES - Legendary (2 items, $50,000+)
		{
			Name:        "Rare Collectible Watch",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityLegendary,
			Price:       85000.00,
			ImageURL:    "/images/accessories/collectible-watch.jpg",
			Description: "Limited edition timepiece from prestigious watchmaker",
		},
		{
			Name:        "Diamond Necklace Set",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityLegendary,
			Price:       120000.00,
			ImageURL:    "/images/accessories/diamond-set.jpg",
			Description: "Exquisite diamond necklace and earring set",
		},

		// CARS - Common (2 items, $1,000-$5,000)
		{
			Name:        "Old Sedan",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       1000.00,
			ImageURL:    "/images/cars/old-sedan.jpg",
			Description: "Reliable but worn sedan",
		},
		{
			Name:        "Compact Car",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       5000.00,
			ImageURL:    "/images/cars/compact-car.jpg",
			Description: "Small and fuel-efficient",
		},

		// CARS - Rare (4 items, $12,000-$45,000)
		{
			Name:        "Family Sedan",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       12000.00,
			ImageURL:    "/images/cars/family-sedan.jpg",
			Description: "Spacious family car",
		},
		{
			Name:        "Used SUV",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       18000.00,
			ImageURL:    "/images/cars/used-suv.jpg",
			Description: "Pre-owned SUV in good condition",
		},
		{
			Name:        "New SUV",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       35000.00,
			ImageURL:    "/images/cars/new-suv.jpg",
			Description: "Brand new SUV",
		},
		{
			Name:        "Sports Coupe",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       45000.00,
			ImageURL:    "/images/cars/sports-coupe.jpg",
			Description: "Fast and stylish coupe",
		},

		// CARS - Epic (3 items, $55,000-$95,000)
		{
			Name:        "Electric Car",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       55000.00,
			ImageURL:    "/images/cars/electric-car.jpg",
			Description: "Modern electric vehicle",
		},
		{
			Name:        "Luxury Sedan",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       60000.00,
			ImageURL:    "/images/cars/luxury-sedan.jpg",
			Description: "High-end luxury sedan",
		},
		{
			Name:        "Tesla Model S",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       95000.00,
			ImageURL:    "/images/cars/tesla-model-s.jpg",
			Description: "Premium electric luxury sedan",
		},

		// CARS - Legendary (4 items, $110,000-$500,000)
		{
			Name:        "Mercedes S-Class",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       110000.00,
			ImageURL:    "/images/cars/mercedes-s.jpg",
			Description: "Ultimate luxury sedan",
		},
		{
			Name:        "Porsche 911",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       125000.00,
			ImageURL:    "/images/cars/porsche-911.jpg",
			Description: "Iconic sports car",
		},
		{
			Name:        "Ferrari F8",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       280000.00,
			ImageURL:    "/images/cars/ferrari-f8.jpg",
			Description: "Italian supercar",
		},
		{
			Name:        "Lamborghini Aventador",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       500000.00,
			ImageURL:    "/images/cars/lamborghini.jpg",
			Description: "Legendary Italian supercar",
		},

		// HOUSES - Common (2 items, $2,000-$5,000)
		{
			Name:        "Studio Apartment",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       5000.00,
			ImageURL:    "/images/houses/studio-apartment.jpg",
			Description: "Cozy studio apartment",
		},
		{
			Name:        "Small Apartment",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       15000.00,
			ImageURL:    "/images/houses/small-apartment.jpg",
			Description: "One-bedroom apartment",
		},

		// HOUSES - Rare (4 items, $35,000-$120,000)
		{
			Name:        "Suburban House",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       35000.00,
			ImageURL:    "/images/houses/suburban-house.jpg",
			Description: "Nice house in the suburbs",
		},
		{
			Name:        "City Condo",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       75000.00,
			ImageURL:    "/images/houses/city-condo.jpg",
			Description: "Modern downtown condo",
		},
		{
			Name:        "Family Home",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       120000.00,
			ImageURL:    "/images/houses/family-home.jpg",
			Description: "Spacious family home",
		},
		{
			Name:        "Lake House",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       95000.00,
			ImageURL:    "/images/houses/lake-house.jpg",
			Description: "Peaceful lakeside retreat",
		},

		// HOUSES - Epic (2 items, $200,000-$500,000)
		{
			Name:        "Beach House",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityEpic,
			Price:       200000.00,
			ImageURL:    "/images/houses/beach-house.jpg",
			Description: "Beautiful beachfront property",
		},
		{
			Name:        "Luxury Penthouse",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityEpic,
			Price:       500000.00,
			ImageURL:    "/images/houses/penthouse.jpg",
			Description: "Top-floor luxury penthouse",
		},

		// HOUSES - Legendary (3 items, $750,000-$1,000,000)
		{
			Name:        "Modern Mansion",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
			Price:       750000.00,
			ImageURL:    "/images/houses/mansion.jpg",
			Description: "Stunning modern mansion",
		},
		{
			Name:        "Private Estate",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
			Price:       1000000.00,
			ImageURL:    "/images/houses/estate.jpg",
			Description: "Exclusive private estate",
		},
		{
			Name:        "Island Villa",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
			Price:       1000000.00,
			ImageURL:    "/images/houses/island-villa.jpg",
			Description: "Private island villa paradise",
		},

		// ========== MEME & WEIRD ITEMS ==========

		// CLOTHING - Meme Items
		{
			Name:        "Tinfoil Hat",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       99.00,
			ImageURL:    "/images/clothing/tinfoil-hat.jpg",
			Description: "Protects you from mind control and alien signals",
		},
		{
			Name:        "Banana Suit",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       2500.00,
			ImageURL:    "/images/clothing/banana-suit.jpg",
			Description: "Go bananas! Full body banana costume",
		},
		{
			Name:        "Dinosaur Pajamas",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       1800.00,
			ImageURL:    "/images/clothing/dino-pajamas.jpg",
			Description: "Rawr! Cozy dino onesie for maximum comfort",
		},
		{
			Name:        "Socks with Sandals",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       150.00,
			ImageURL:    "/images/clothing/socks-sandals.jpg",
			Description: "The ultimate fashion faux pas. Dad-approved!",
		},
		{
			Name:        "Ugly Christmas Sweater",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       999.00,
			ImageURL:    "/images/clothing/ugly-sweater.jpg",
			Description: "So ugly it's beautiful. Perfect for parties!",
		},
		{
			Name:        "Chicken Suit",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       3200.00,
			ImageURL:    "/images/clothing/chicken-suit.jpg",
			Description: "Don't be chicken! Embrace your inner poultry",
		},
		{
			Name:        "Cat Meme T-Shirt",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityCommon,
			Price:       420.00,
			ImageURL:    "/images/clothing/cat-tshirt.jpg",
			Description: "I Can Has Cheezburger? Classic internet culture",
		},
		{
			Name:        "Unicorn Onesie",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       2800.00,
			ImageURL:    "/images/clothing/onesie.jpg",
			Description: "Magical unicorn pajamas with rainbow tail",
		},
		{
			Name:        "Inflatable T-Rex Costume",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityEpic,
			Price:       8500.00,
			ImageURL:    "/images/clothing/inflatable-dinosaur.jpg",
			Description: "Battery-powered inflatable dinosaur suit. Pure chaos!",
		},
		{
			Name:        "Pickle Costume",
			Type:        model.ItemTypeClothing,
			Rarity:      model.ItemRarityRare,
			Price:       3500.00,
			ImageURL:    "/images/clothing/pickle-costume.jpg",
			Description: "I'm Pickle Rick! The funniest thing you've ever seen",
		},

		// ACCESSORIES - Meme Items
		{
			Name:        "Rubber Duck",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       69.00,
			ImageURL:    "/images/accessories/rubber-duck.jpg",
			Description: "Debugging companion. Quack quack!",
		},
		{
			Name:        "Fake Mustache",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       199.00,
			ImageURL:    "/images/accessories/fake-mustache.jpg",
			Description: "Instant disguise! Look sophisticated instantly",
		},
		{
			Name:        "Monocle",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       4200.00,
			ImageURL:    "/images/accessories/monocle.jpg",
			Description: "Indeed! Quite fancy, old chap",
		},
		{
			Name:        "Power Ring",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityEpic,
			Price:       9999.00,
			ImageURL:    "/images/accessories/power-ring.jpg",
			Description: "One ring to rule them all (not really)",
		},
		{
			Name:        "Magic Wand",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       6969.00,
			ImageURL:    "/images/accessories/magic-wand.jpg",
			Description: "Wingardium Leviosa! May or may not actually work",
		},
		{
			Name:        "Glowing LED Glasses",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityRare,
			Price:       3333.00,
			ImageURL:    "/images/accessories/glowing-glasses.jpg",
			Description: "Cyberpunk vibes. Perfect for raves!",
		},
		{
			Name:        "Potato Clock",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       777.00,
			ImageURL:    "/images/accessories/potato-clock.jpg",
			Description: "Powered by a potato. Science!",
		},
		{
			Name:        "Finger Hands",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       555.00,
			ImageURL:    "/images/accessories/finger-hands.jpg",
			Description: "Tiny hands for your fingers. Cursed!",
		},
		{
			Name:        "Googly Eyes Glasses",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       299.00,
			ImageURL:    "/images/accessories/googly-eyes.jpg",
			Description: "Make everyone laugh with these silly specs",
		},
		{
			Name:        "Unicorn Horn Headband",
			Type:        model.ItemTypeAccessories,
			Rarity:      model.ItemRarityCommon,
			Price:       888.00,
			ImageURL:    "/images/accessories/unicorn-horn.jpg",
			Description: "Become a magical unicorn instantly",
		},

		// CARS - Meme Items
		{
			Name:        "Oscar Mayer Wienermobile",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       75000.00,
			ImageURL:    "/images/cars/hotdog-car.jpg",
			Description: "The legendary hot dog car! I wish I were an Oscar Mayer Wiener",
		},
		{
			Name:        "Flying Carpet",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       42000.00,
			ImageURL:    "/images/cars/magic-carpet.jpg",
			Description: "A whole new world! (Warning: May not actually fly)",
		},
		{
			Name:        "Racing Shopping Cart",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       666.00,
			ImageURL:    "/images/cars/shopping-cart.jpg",
			Description: "Supermarket sweep champion! Squeaky wheel included",
		},
		{
			Name:        "Segway",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       8999.00,
			ImageURL:    "/images/cars/segway.jpg",
			Description: "Mall cop approved transportation",
		},
		{
			Name:        "Power Wheels Barbie Jeep",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       1999.00,
			ImageURL:    "/images/cars/power-wheels.jpg",
			Description: "Limited to 5mph. Living the dream!",
		},
		{
			Name:        "Ice Cream Truck",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       28000.00,
			ImageURL:    "/images/cars/ice-cream-truck.jpg",
			Description: "Plays music everywhere you go. Kids love it!",
		},
		{
			Name:        "Adult Tricycle",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityCommon,
			Price:       899.00,
			ImageURL:    "/images/cars/tricycle.jpg",
			Description: "Three wheels of pure stability",
		},
		{
			Name:        "Rocket-Powered Skateboard",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityEpic,
			Price:       15000.00,
			ImageURL:    "/images/cars/skateboard.jpg",
			Description: "Back to the Future style! Safety not guaranteed",
		},
		{
			Name:        "Hoverboard",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityRare,
			Price:       6666.00,
			ImageURL:    "/images/cars/hoverboard.jpg",
			Description: "Doesn't actually hover, but it's still cool",
		},
		{
			Name:        "Toy Rocket Ship",
			Type:        model.ItemTypeCar,
			Rarity:      model.ItemRarityLegendary,
			Price:       99999.00,
			ImageURL:    "/images/cars/rocket.jpg",
			Description: "To the moon! Literally!",
		},

		// HOUSES - Meme Items
		{
			Name:        "Cardboard Box Mansion",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       420.00,
			ImageURL:    "/images/houses/cardboard-box.jpg",
			Description: "Refrigerator box deluxe. Peak minimalist living",
		},
		{
			Name:        "Epic Treehouse",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       25000.00,
			ImageURL:    "/images/houses/treehouse.jpg",
			Description: "No parents allowed! Secret password required",
		},
		{
			Name:        "Cozy Cave",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       3000.00,
			ImageURL:    "/images/houses/cave.jpg",
			Description: "Return to monke. Unga bunga!",
		},
		{
			Name:        "Pillow Fort Palace",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       1337.00,
			ImageURL:    "/images/houses/pillow-fort.jpg",
			Description: "Engineered for maximum comfort. Blankets included!",
		},
		{
			Name:        "Arctic Igloo",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       18000.00,
			ImageURL:    "/images/houses/igloo.jpg",
			Description: "Ice ice baby! Natural air conditioning",
		},
		{
			Name:        "Hobbit Hole",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityEpic,
			Price:       88888.00,
			ImageURL:    "/images/houses/hobbit-hole.jpg",
			Description: "In a hole in the ground there lived a hobbit",
		},
		{
			Name:        "UFO Landing Base",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityLegendary,
			Price:       420000.00,
			ImageURL:    "/images/houses/ufo.jpg",
			Description: "Aliens welcome! Free WiFi included",
		},
		{
			Name:        "Luxury Camping Tent",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityCommon,
			Price:       2500.00,
			ImageURL:    "/images/houses/tent.jpg",
			Description: "Glamping at its finest",
		},
		{
			Name:        "Van Down By The River",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityRare,
			Price:       15000.00,
			ImageURL:    "/images/houses/van.jpg",
			Description: "Living the #VanLife dream! Government cheese optional",
		},
		{
			Name:        "Inflatable Bouncy Castle",
			Type:        model.ItemTypeHouse,
			Rarity:      model.ItemRarityEpic,
			Price:       50000.00,
			ImageURL:    "/images/houses/castle.jpg",
			Description: "Your childhood dream home. Bouncing is mandatory!",
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
