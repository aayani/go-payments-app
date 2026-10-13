package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, accounts *AccountHandler, transactions *TransactionHandler) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("/api")
	api.POST("/transactions", transactions.CreateTransaction)
	api.POST("/transactions/init", transactions.InitTransaction)

	admin := api.Group("/admin")
	admin.GET("/accounts", accounts.GetAccounts)
	admin.POST("/accounts", accounts.CreateAccount)
	admin.GET("/transactions/:accountNumber", transactions.GetTransactionsByAccountNumber)
}
