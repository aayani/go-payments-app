package models

import "errors"

var (
	ErrInsufficientBalance       = errors.New("insufficient balance")
	ErrInvalidAccountNumber      = errors.New("invalid account number")
	ErrTransactionUnavailable    = errors.New("transaction is unavailable")
	ErrInvalidTransactionRequest = errors.New("invalid transaction request")
	ErrAccountsUnavailable       = errors.New("source or destination account does not exist")
	ErrInvalidAccountTitle       = errors.New("title is required and must be at most 255 characters")
)
