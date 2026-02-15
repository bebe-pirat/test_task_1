package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test_task/internal/entity"
	"test_task/internal/service"

	"github.com/gorilla/mux"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) CreateSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Неправильный джейсон", err)
		return
	}

	err := h.service.CreateSubscription(ctx, request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка создания подписки", err)
	}
}

func (h *SubscriptionHandler) DeleteSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid json", http.StatusInternalServerError)
		log.Println("Неверный json", err)
		return
	}

	err = h.service.DeleteSubById(ctx, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка удаления подписки", err)
	}
}

func (h *SubscriptionHandler) UpdateSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Println("Неправильный джейсон", err)
		return
	}

	err := h.service.UpdateSubById(ctx, request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка обновления подписки", err)
	}
}

func (h *SubscriptionHandler) GetAllSubsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := h.service.GetAllSubscriptions(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка чтения подписок", err)
	}

	jsonData, err := json.Marshal(subs)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка сериализации", err)
		return
	}

	w.Write(jsonData)
}

func (h *SubscriptionHandler) GetSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid json", http.StatusInternalServerError)
		log.Println("Неверный json", err)
		return
	}

	sub, err := h.service.GetSubscriptionById(ctx, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка чтения подписок", err)
		return
	}

	jsonData, err := json.Marshal(sub)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Ошибка сериализации", err)
		return
	}

	w.Write(jsonData)
}
