package main

import (
	"log/slog"
	"net/http"
	"os"
	"test_task/internal/database"
	"test_task/internal/handler"
	"test_task/internal/repository"
	"test_task/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using system environment variables")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := database.InitDB(connStr)
	if err != nil {
		slog.Error("Ошибка инициализации БД: ", "error", err)
		return
	}
	defer database.CloseDB(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
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
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	slog.Info("server started" + port)

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server didn't started", "error", err)
		return
	}
}
