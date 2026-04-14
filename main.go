package main

import (
	"log"
	"net/http"
	"task-manager/internal/database"
	"task-manager/internal/routes"
	"task-manager/utils/logger"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logger.Init()

	logger.Info("Starting task-manager API", "version", "1.0.0")

	db := database.Connect()
	logger.Info("Database connected successfully")

	r := routes.SetupRouter(db)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	logger.Info("Server listening", "port", "8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server error", "error", err.Error())
	}
}
