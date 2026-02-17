package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"test_task/internal/entity"
	"test_task/internal/service"

	"github.com/google/uuid"
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
		slog.Error("Неправильный JSON", "error", err)
		return
	}
	defer r.Body.Close()

	_, err := h.service.CreateSubscription(ctx, request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка создания подписки", "error", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SubscriptionHandler) DeleteSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		slog.Error("Неверный json", "error", err)
		return
	}

	err = h.service.DeleteSubById(ctx, id)
	if err == sql.ErrNoRows {
		http.Error(w, "Row not found", http.StatusNotFound)
		slog.Info("Не найдена запись для удаления", "error", sql.ErrNoRows)
		return
	}

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка удаления подписки", "error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) UpdateSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		slog.Error("Неправильный джейсон", "error", err)
		return
	}
	defer r.Body.Close()

	err := h.service.UpdateSubById(ctx, request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка обновления подписки", "error", err)
		return
	}
}

func (h *SubscriptionHandler) GetAllSubsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := h.service.GetAllSubscriptions(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка чтения подписок", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(subs); err != nil {
		slog.Error("Ошибка сериализации", "error", err)
	}
}

func (h *SubscriptionHandler) GetSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		slog.Error("Неверный json", "error", err)
		return
	}

	sub, err := h.service.GetSubscriptionById(ctx, id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка чтения подписок", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(sub); err != nil {
		slog.Error("Ошибка сериализации", "error", err)
	}
}

func (h *SubscriptionHandler) GetTotalCostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()

	serviceName := query.Get("service_name")
	userIDStr := query.Get("user_id")
	fromDate := query.Get("from_date")
	toDateStr := query.Get("to_date")
	var toDate *string = nil
	if toDateStr != "" {
		toDate = &toDateStr
	}

	var userID uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			slog.Error("ошибка парсинга id", "error", err)
			return
		}
		userID = parsed
	}

	total, err := h.service.GetTotalCost(ctx, userID, serviceName, fromDate, toDate)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		slog.Error("ошибка получения фильтрации", "error", err)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{
		"total": total,
	})
}
