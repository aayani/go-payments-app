package controllers

import (
	"net/http"

	"github.com/aayani/go-payments-app/internal/models"
	"github.com/aayani/go-payments-app/internal/services"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	Payments *services.PaymentService
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var request models.CreateAccountPayload
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": models.ErrInvalidAccountTitle.Error()})
		return
	}
	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.Payments.CreateAccount(c.Request.Context(), request.Title)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}
	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.Payments.ListAccounts(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list accounts"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}
