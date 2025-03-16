package main

import (
	"context"
	"net/http"

	"github.com/grnsv/gophermart/internal/api/handlers"
	"github.com/grnsv/gophermart/internal/api/router"
	"github.com/grnsv/gophermart/internal/config"
	"github.com/grnsv/gophermart/internal/logger"
	"github.com/grnsv/gophermart/internal/services"
	"github.com/grnsv/gophermart/internal/storage"
)

func main() {
	log := logger.New()
	defer log.Sync()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Infof("Starting server with config: %v", cfg)
	store, err := storage.New(context.Background(), cfg.DatabaseURI)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer store.Close()

	jwtService := services.NewJWTService(cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(log, services.NewUserService(store), jwtService)
	orderHandler := handlers.NewOrderHandler(log, services.NewOrderService(store), services.NewLuhnService())
	router := router.NewRouter(log, userHandler, orderHandler, jwtService)
	if err := http.ListenAndServe(cfg.RunAddress, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
