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

type Application struct {
	Logger         logger.Logger
	Config         *config.Config
	Storage        storage.Storage
	JWTService     services.JWTService
	OrderService   services.OrderService
	AccrualService services.AccrualService
	UserService    services.UserService
	Validator      services.Validator
	Router         http.Handler
}

func NewApplication(ctx context.Context) *Application {
	var app Application

	app.initLogger()
	app.initConfig()
	app.initStorage(ctx)
	app.initServices()
	app.initHandlers()

	return &app
}

func (app *Application) initLogger() {
	app.Logger = logger.New()
}

func (app *Application) initConfig() {
	var err error
	app.Config, err = config.New()
	if err != nil {
		app.Logger.Fatalf("Failed to load config: %v", err)
	}
}

func (app *Application) initStorage(ctx context.Context) {
	var err error
	app.Storage, err = storage.New(ctx, app.Config.DatabaseURI)
	if err != nil {
		app.Logger.Fatalf("Failed to create storage: %v", err)
	}
}

func (app *Application) initServices() {
	app.JWTService = services.NewJWTService(app.Config.JWTSecret)
	app.AccrualService = services.NewAccrualService(app.Config.AccrualSystemAddress)
	app.OrderService = services.NewOrderService(app.Logger, app.Storage, app.AccrualService)
	app.UserService = services.NewUserService(app.Storage)
	app.Validator = services.NewLuhnService()
}

func (app *Application) initHandlers() {
	authHandler := handlers.NewAuthHandler(app.Logger, app.UserService, app.JWTService)
	protectedHandler := handlers.NewProtectedHandler(app.Logger, app.OrderService, app.Validator, app.Storage)
	app.Router = router.NewRouter(app.Logger, authHandler, protectedHandler, app.JWTService)
}

func (app *Application) Run() {
	app.Logger.Infof("Starting server with config: %v", app.Config)
	if err := http.ListenAndServe(app.Config.RunAddress, app.Router); err != nil {
		app.Logger.Fatalf("Server failed: %v", err)
	}
}

func (app *Application) Close() {
	app.Storage.Close()
	app.Logger.Sync()
}

func main() {
	app := NewApplication(context.Background())
	defer app.Close()

	app.Run()
}
