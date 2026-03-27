package main

import (
	"log"
	"os"
	"smart-inventory/internal/handler"
	"smart-inventory/internal/repository"
	"smart-inventory/internal/service"
	"smart-inventory/pkg/database"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := "postgres://" +
		os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASSWORD") + "@" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + "/" +
		os.Getenv("DB_NAME") +
		"?search_path=" + os.Getenv("DB_SCHEMA") +
		"&sslmode=" + os.Getenv("DB_SSLMODE")

	db, err := database.NewPostgres(dbURL)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	invRepo := repository.NewInventoryRepository(db)

	invService := service.NewInventoryService(invRepo, db)
	invHandler := handler.NewInventoryHandler(invService)

	r.POST("/inventory/adjust", invHandler.AdjustStock)
	r.GET("/inventory", invHandler.GetInventory)

	// STOCK IN
	stockRepo := repository.NewStockInRepository()
	stockService := service.NewStockInService(db, stockRepo, invRepo)
	stockHandler := handler.NewStockInHandler(stockService)

	r.POST("/stock-in", stockHandler.Create)
	r.PUT("/stock-in/:id/status", stockHandler.UpdateStatus)
	r.GET("/stock-in", stockHandler.GetAll)

	stockOutRepo := repository.NewStockOutRepository()
	stockOutService := service.NewStockOutService(db, stockOutRepo, invRepo)
	stockOutHandler := handler.NewStockOutHandler(stockOutService)

	r.POST("/stock-out", stockOutHandler.Create)
	r.POST("/stock-out/:id/in-progress", stockOutHandler.InProgress)
	r.POST("/stock-out/:id/complete", stockOutHandler.Complete)
	r.POST("/stock-out/:id/cancel", stockOutHandler.Cancel)
	r.GET("/stock-out", stockOutHandler.GetAll)

	r.Run(":9090")
}
