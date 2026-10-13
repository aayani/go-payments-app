package postgres

import (
	"fmt"

	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OpenDB opens a GORM connection to PostgreSQL.
func OpenDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}
	return db, nil
}

// CloseDB closes the underlying PostgreSQL connection.
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close PostgreSQL connection: %w", err)
	}
	return nil
}
