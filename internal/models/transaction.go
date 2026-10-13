package models

import (
	"strconv"
	"strings"
	"time"
)

type Transaction struct {
	ID            string     `gorm:"type:uuid;primaryKey;not null"`
	ToAccountID   string     `gorm:"type:uuid;not null" json:"-"`
	FromAccountID string     `gorm:"type:uuid;not null" json:"-"`
	Amount        float64    `gorm:"type:decimal(10,2);not null"`
	CreatedAt     *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt     *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	ToAccount     Account    `gorm:"foreignKey:ToAccountID;references:ID;constraint:transactions_to_account_id_fkey,OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	FromAccount   Account    `gorm:"foreignKey:FromAccountID;references:ID;constraint:transactions_from_account_id_fkey,OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

type InitTransactionPayload struct {
	Source      string  `json:"source" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

type CreateTransactionPayload struct {
	TransactionID string `json:"transactionId" binding:"required,uuid"`
}

func (payload *InitTransactionPayload) Validate() error {
	_, fraction, _ := strings.Cut(strconv.FormatFloat(payload.Amount, 'f', -1, 64), ".")
	if !AccountNumberPattern.MatchString(payload.Source) || !AccountNumberPattern.MatchString(payload.Destination) || payload.Source == payload.Destination ||
		!(payload.Amount > 0 && payload.Amount < 100000000) || len(fraction) > 2 {
		return ErrInvalidTransactionRequest
	}
	return nil
}
