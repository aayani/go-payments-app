package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aayani/go-payments-app/app"
	"github.com/aayani/go-payments-app/internal/controllers"
	"github.com/aayani/go-payments-app/internal/postgres"
	"github.com/aayani/go-payments-app/internal/postgres/migrations"
	"github.com/aayani/go-payments-app/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("failed to load .env; check its syntax and permissions")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	db, err := postgres.OpenDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.CloseDB(db)
	if err := migrations.Migrate(context.Background(), db); err != nil {
		log.Fatal("migration failed: ", err)
	}
	log.Print("migration complete")

	if err := migrations.Seed(context.Background(), db); err != nil {
		log.Fatal("seed failed: ", err)
	}
	log.Print("seed complete")

	redisOptions, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("REDIS_URL is required and must be a valid Redis URL")
	}
	redisClient := redis.NewClient(redisOptions)
	defer redisClient.Close()

	router := gin.Default()
	repository := &postgres.Repository{DB: db}
	accountHandler := &controllers.AccountHandler{Payments: &services.PaymentService{Repository: repository}}
	transactionHandler := &controllers.TransactionHandler{Transactions: &services.TransactionService{Repository: repository, Redis: redisClient}}
	controllers.RegisterRoutes(router, accountHandler, transactionHandler)

	router.StaticFS("/docs", http.FS(app.Files))

	if err := router.Run(); err != nil {
		log.Fatal(err)
	}
}
