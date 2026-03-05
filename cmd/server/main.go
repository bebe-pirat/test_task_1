package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test_task/internal/database"
	"test_task/internal/handler"
	"test_task/internal/repository"
	"test_task/internal/service"
	"time"

	_ "test_task/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Package main API
//
// @title Subscription Service API
// @version 1.0
// @description API for managing subscriptions
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.email support@example.com
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @tag.name subscriptions
// @tag.description Subscription management endpoints

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using system environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := database.InitDB(connStr)
	if err != nil {
		slog.Error("database is not initialized", "error", err)
		return
	}
	defer database.CloseDB(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	subRepo := repository.NewSubscriptionRepository(db)
	subService := service.NewSubscriptionService(subRepo)
	subHandler := handler.NewSubscriptionHandler(subService)

	router := mux.NewRouter()

	router.HandleFunc("/subscriptions", subHandler.CreateSubHandler).Methods("POST")
	router.HandleFunc("/subscriptions", subHandler.GetAllSubsHandler).Methods("GET")
	router.HandleFunc("/subscriptions/total", subHandler.GetTotalCostHandler).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", subHandler.GetSubHandler).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", subHandler.DeleteSubHandler).Methods("DELETE")
	router.HandleFunc("/subscriptions/{id}", subHandler.UpdateSubHandler).Methods("PUT")
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("./swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		slog.Info("server started", "port", port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server didn't started", "error", err)
			return
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	slog.Info("server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exiting")
}
