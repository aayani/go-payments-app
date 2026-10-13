package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aayani/go-payments-app/internal/models"
	"github.com/aayani/go-payments-app/internal/postgres"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TransactionService struct {
	Redis      *redis.Client
	Repository *postgres.Repository
}

func (service *TransactionService) GenerateTransactionId(ctx context.Context, payload *models.InitTransactionPayload) (string, error) {
	accounts, err := service.Repository.GetAccountsByAccountNumbers(ctx, []string{payload.Source, payload.Destination})
	if err != nil {
		return "", err
	}
	if accounts == nil {
		return "", models.ErrAccountsUnavailable
	}

	for _, account := range accounts {
		if account.AccountNumber == payload.Source {
			if payload.Amount > account.Balance {
				return "", models.ErrInsufficientBalance
			}
			break
		}
	}

	id := uuid.NewString()

	value, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return id, service.Redis.Set(ctx, id, string(value), 5*time.Minute).Err()
}

func (service *TransactionService) CreateTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	value, err := service.Redis.GetDel(ctx, transactionID).Result()
	if errors.Is(err, redis.Nil) {
		return nil, models.ErrTransactionUnavailable
	}
	if err != nil {
		return nil, err
	}
	var payload models.InitTransactionPayload
	if err := json.Unmarshal([]byte(value), &payload); err != nil {
		return nil, err
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	transaction := &models.Transaction{
		ID:          transactionID,
		Amount:      payload.Amount,
		FromAccount: models.Account{AccountNumber: payload.Source},
		ToAccount:   models.Account{AccountNumber: payload.Destination},
	}
	if err := service.Repository.CreateTransaction(ctx, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *TransactionService) ListTransactions(ctx context.Context, accountNumber string) ([]models.Transaction, error) {
	return s.Repository.ListTransactions(ctx, accountNumber)
}
