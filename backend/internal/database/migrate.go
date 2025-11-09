package database

import (
	"fmt"
	"log"

	"github.com/smoreg/freezino/backend/internal/model"
)

// Migrate runs auto-migration for all models
func Migrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Running database migrations...")

	// Auto-migrate all models
	err := DB.AutoMigrate(
		&model.User{},
		&model.Transaction{},
		&model.Item{},
		&model.UserItem{},
		&model.WorkSession{},
		&model.GameSession{},
		&model.ContactMessage{},
		&model.RouletteResult{},
		&model.UserStatus{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")

	// Seed initial data (items) as part of migration
	log.Println("Checking for initial data...")
	if err := seedInitialData(); err != nil {
		log.Printf("Warning: Failed to seed initial data: %v", err)
		// Don't fail the migration if seeding fails
	}

	return nil
}

// DropAllTables drops all tables (use with caution!)
func DropAllTables() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Dropping all tables...")

	err := DB.Migrator().DropTable(
		&model.UserStatus{},
		&model.RouletteResult{},
		&model.ContactMessage{},
		&model.GameSession{},
		&model.WorkSession{},
		&model.UserItem{},
		&model.Transaction{},
		&model.Item{},
		&model.User{},
	)

	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	log.Println("All tables dropped successfully")
	return nil
}

// ResetDatabase drops and recreates all tables
func ResetDatabase() error {
	if err := DropAllTables(); err != nil {
		return err
	}

	return Migrate()
}

// seedInitialData seeds the database with initial items
// This is called as part of migrations
func seedInitialData() error {
	// Check if items already exist
	var count int64
	DB.Model(&model.Item{}).Count(&count)

	if count > 0 {
		log.Printf("Initial data already exists (%d items), skipping seed...", count)
		return nil
	}

	log.Println("Seeding initial items...")

	// Call the seedItems function from seed.go
	if err := seedItems(); err != nil {
		return fmt.Errorf("failed to seed items: %w", err)
	}

	log.Println("Initial data seeded successfully")
	return nil
}
