package handler

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
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

// CreateSubHandler godoc
//
// @Summary Create a new subscription
// @Tags subscriptions
// @Accept json
// @Produce application/json
// @Param subscription body entity.Subscription true "Subscription data"
// @Success 201
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions [post]
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

// DeleteSubHandler godoc
//
// @Summary Delete a subscription by id
// @Description Delete an existing subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce application/json
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
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

// UpdateSubHandler godoc
// @Summary Update a subscription by id
// @Description Update an existing subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body entity.Subscription true "Subscription data"
// @Success 204
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request entity.Subscription

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		slog.Error("Неправильный JSON", "error", err)
		return
	}
	defer r.Body.Close()

	err := h.service.UpdateSubById(ctx, request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка обновления подписки", "error", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAllSubsHandler godoc
// @Summary Get all subscriptions
// @Description Show all existing subscriptions
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} entity.Subscription "List of all subscriptions"
// @Failure 500 {string} string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAllSubsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := h.service.GetAllSubscriptions(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		slog.Error("Ошибка чтения подписок", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(subs); err != nil {
		slog.Error("Ошибка сериализации", "error", err)
	}
}

// GetSubHandler godoc
// @Summary Get a subscription by id
// @Description Show an existing subscription by id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} entity.Subscription "Subscription found"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(sub); err != nil {
		slog.Error("Ошибка сериализации", "error", err)
	}
}

// GetTotalCostHandler godoc
// @Summary Get total cost of subscriptions
// @Description Calculate total cost of subscriptions with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by user ID" Format(uuid)
// @Param service_name query string false "Filter by service name"
// @Param from_date query string false "Filter by start date" Format(date)
// @Param to_date query string false "Filter by end date" Format(date)
// @Success 200 {object} map[string]int "Total cost"
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /subscriptions/total [get]
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]int{
		"total": total,
	}); err != nil {
		slog.Error("ошибка сериализации", "error", err)
	}
}
