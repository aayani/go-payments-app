package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/aayani/go-payments-app/internal/models"
	"github.com/aayani/go-payments-app/internal/postgres"
	"github.com/google/uuid"
)

type PaymentService struct {
	Repository *postgres.Repository
}

func (s *PaymentService) CreateAccount(ctx context.Context, title string) (models.Account, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(100_000_000_000))
	if err != nil {
		return models.Account{}, err
	}
	account := models.Account{
		ID:            uuid.NewString(),
		AccountNumber: fmt.Sprintf("%04d-%07d", number.Int64()/10_000_000, number.Int64()%10_000_000),
		Title:         title,
	}
	if err := s.Repository.CreateAccount(ctx, &account); err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func (s *PaymentService) ListAccounts(ctx context.Context) ([]models.Account, error) {
	return s.Repository.ListAccounts(ctx)
}
