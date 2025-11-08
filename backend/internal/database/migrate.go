package database

import (
	"fmt"
	"log"

	"freezino/backend/internal/model"
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
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// DropAllTables drops all tables (use with caution!)
func DropAllTables() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	log.Println("Dropping all tables...")

	err := DB.Migrator().DropTable(
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
