package migrations

import (
	"context"

	"github.com/aayani/go-payments-app/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Seed(ctx context.Context, db *gorm.DB) error {
	accounts := []models.Account{
		{ID: "00000000-0000-4000-8000-000000000001", AccountNumber: "0001-0000001", Title: "Alice", Balance: 1000},
		{ID: "00000000-0000-4000-8000-000000000002", AccountNumber: "0001-0000002", Title: "Bob", Balance: 2000},
		{ID: "00000000-0000-4000-8000-000000000003", AccountNumber: "0001-0000003", Title: "Carol", Balance: 3000},
		{ID: "00000000-0000-4000-8000-000000000004", AccountNumber: "0001-0000004", Title: "Dave", Balance: 4000},
		{ID: "00000000-0000-4000-8000-000000000005", AccountNumber: "0001-0000005", Title: "Eve", Balance: 5000},
	}
	insertOnly := clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoNothing: true}
	return db.WithContext(ctx).Clauses(insertOnly).Create(&accounts).Error
}
