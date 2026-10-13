package postgres

import (
	"context"
	"strconv"

	"github.com/aayani/go-payments-app/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) ListAccounts(ctx context.Context) ([]models.Account, error) {
	accounts := []models.Account{}
	err := r.DB.WithContext(ctx).Order("id").Find(&accounts).Error
	return accounts, err
}

func (r *Repository) GetAccountsByAccountNumbers(ctx context.Context, accountNumbers []string) ([]models.Account, error) {
	var accounts []models.Account
	err := r.DB.WithContext(ctx).Where("account_number IN ?", accountNumbers).Find(&accounts).Error
	if err != nil || len(accounts) != len(accountNumbers) {
		return nil, err
	}
	return accounts, nil
}

func (r *Repository) CreateAccount(ctx context.Context, account *models.Account) error {
	return r.DB.WithContext(ctx).Create(account).Error
}

func (r *Repository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	if transaction.ID == "" {
		transaction.ID = uuid.NewString()
	}
	amount := strconv.FormatFloat(transaction.Amount, 'f', -1, 64)
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var accounts []models.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("account_number IN ?", []string{transaction.FromAccount.AccountNumber, transaction.ToAccount.AccountNumber}).
			Order("id").Find(&accounts).Error; err != nil {
			return err
		}
		if len(accounts) != 2 {
			return gorm.ErrRecordNotFound
		}
		for _, account := range accounts {
			if account.AccountNumber == transaction.FromAccount.AccountNumber {
				transaction.FromAccountID = account.ID
			} else {
				transaction.ToAccountID = account.ID
			}
		}
		for id, change := range map[string]string{transaction.FromAccountID: "-" + amount, transaction.ToAccountID: amount} {
			query := tx.Model(&models.Account{}).Where("id = ?", id)
			if id == transaction.FromAccountID {
				query = query.Where("balance >= CAST(? AS numeric)", amount)
			}
			result := query.Update("balance", gorm.Expr("balance + CAST(? AS numeric)", change))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return models.ErrInsufficientBalance
			}
		}
		if err := tx.Omit("ToAccount", "FromAccount").Create(transaction).Error; err != nil {
			return err
		}
		return tx.Preload("ToAccount").Preload("FromAccount").
			First(transaction, "id = ?", transaction.ID).Error
	})
}

func (r *Repository) ListTransactions(ctx context.Context, accountNumber string) ([]models.Transaction, error) {
	account := r.DB.WithContext(ctx).Model(&models.Account{}).Select("id").Where("account_number = ?", accountNumber)
	transactions := []models.Transaction{}
	err := r.DB.WithContext(ctx).Preload("ToAccount").Preload("FromAccount").
		Where("from_account_id IN (?) OR to_account_id IN (?)", account, account).
		Order("id").Find(&transactions).Error
	return transactions, err
}
