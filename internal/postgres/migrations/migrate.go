package migrations

import (
	"context"

	"github.com/aayani/go-payments-app/internal/models"
	"gorm.io/gorm"
)

func Migrate(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.AutoMigrate(&models.Account{}, &models.Transaction{})
	})
}
