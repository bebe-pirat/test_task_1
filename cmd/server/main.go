package main

import (
	"log"
	"net/http"
	"os"
	"test_task/internal/database"
	"test_task/internal/handler"
	"test_task/internal/repository"
	"test_task/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	var db database.Database
	err := db.InitDB()
	if err != nil {
		log.Printf("Ошибка инициализации БД: %v\n", err)
		return
	}
	defer db.CloseDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	subRepo := repository.NewSubscriptionRepository(db.DB)
	subService := service.NewSubscriptionService(subRepo)
	subHandler := handler.NewSubscriptionHandler(subService)

	router := mux.NewRouter()

	router.HandleFunc("/subscriptions", subHandler.GetAllSubsHandler).Methods("GET")
	router.HandleFunc("/subscriptions/id", subHandler.GetSubHandler).Methods("GET")
	http.Handle("/", router)

	log.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
