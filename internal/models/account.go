package models

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var AccountNumberPattern = regexp.MustCompile(`^[0-9]{4}-[0-9]{7}$`)

type Account struct {
	ID            string     `gorm:"type:uuid;primaryKey;not null"`
	AccountNumber string     `gorm:"type:varchar(12);not null;uniqueIndex;check:account_number_format,account_number ~ '^[0-9]{4}-[0-9]{7}$';<-:create"`
	Title         string     `gorm:"type:varchar(255);not null"`
	Balance       float64    `gorm:"type:decimal(10,2);not null"`
	CreatedAt     *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt     *time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type CreateAccountPayload struct {
	Title string `json:"title"`
}

func (payload *CreateAccountPayload) Validate() error {
	payload.Title = strings.TrimSpace(payload.Title)
	if payload.Title == "" || utf8.RuneCountInString(payload.Title) > 255 {
		return ErrInvalidAccountTitle
	}
	return nil
}
