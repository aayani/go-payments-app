package controllers

import (
	"errors"
	"net/http"

	"github.com/aayani/go-payments-app/internal/models"
	"github.com/aayani/go-payments-app/internal/services"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	Transactions *services.TransactionService
}

func (handler *TransactionHandler) InitTransaction(ctx *gin.Context) {
	var request models.InitTransactionPayload
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidTransactionRequest.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidTransactionRequest.Error()})
		return
	}
	id, err := handler.Transactions.GenerateTransactionId(ctx.Request.Context(), &request)
	if err != nil {
		ctx.Error(err)
		status := http.StatusInternalServerError
		if errors.Is(err, models.ErrAccountsUnavailable) || errors.Is(err, models.ErrInsufficientBalance) {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": "failed to initialize transaction"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"transactionId": id})
}

func (h *TransactionHandler) CreateTransaction(ctx *gin.Context) {
	var request models.CreateTransactionPayload
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidTransactionRequest.Error()})
		return
	}
	transaction, err := h.Transactions.CreateTransaction(ctx.Request.Context(), request.TransactionID)
	if err != nil {
		ctx.Error(err)
		status := http.StatusInternalServerError
		if errors.Is(err, models.ErrTransactionUnavailable) || errors.Is(err, models.ErrInsufficientBalance) {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{"error": "transaction could not be completed"})
		return
	}
	ctx.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransactionsByAccountNumber(c *gin.Context) {
	accountNumber := c.Param("accountNumber")
	if !models.AccountNumberPattern.MatchString(accountNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidAccountNumber.Error()})
		return
	}
	transactions, err := h.Transactions.ListTransactions(c.Request.Context(), accountNumber)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list transactions"})
		return
	}
	c.JSON(http.StatusOK, transactions)
}
