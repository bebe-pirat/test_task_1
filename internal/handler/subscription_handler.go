package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test_task/internal/entity"
	"test_task/internal/repository"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	subRepo *repository.SubscriptionRepository
}

func (h *SubscriptionHandler) CreateSubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Метод не разрешен")
		return
	}

	var request entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Неправильный джейсон")
		return
	}

	err := h.subRepo.CreateSubscription(request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка создания подписки")
	}
}

func (h *SubscriptionHandler) DeleteSubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Метод не разрешен")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid json", http.StatusInternalServerError)
		log.Println("Неверный json")
		return
	}

	err = h.subRepo.DeleteSubById(id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка удаления подписки")
	}
}

func (h *SubscriptionHandler) UpdateSubHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Метод не разрешен")
		return
	}

	var request entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Неправильный джейсон")
		return
	}

	err := h.subRepo.UpdateSubById(request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка обновления подписки")
	}
}

func (h *SubscriptionHandler) ReadSubscriptions(w http.ResponseWriter, r http.Request) {

}
