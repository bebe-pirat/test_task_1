package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"test_task/internal/entity"
	"test_task/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, e entity.Subscription) (int, error) {
	if e.ServiceName == "" {
		return 0, errors.New("service name is required")
	}

	if e.Price < 0 {
		return 0, errors.New("price should be non-negative")
	}

	err := isDateValid(&e.StartDate, e.EndDate)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateSubscription(ctx, e)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SubscriptionService) GetSubscriptionById(ctx context.Context, id int) (*entity.Subscription, error) {
	if id <= 0 {
		return nil, errors.New("subscription id is required")
	}

	return s.repo.GetSubscriptionById(ctx, id)
}

func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	return s.repo.GetAllSubscriptions(ctx)
}

func (s *SubscriptionService) DeleteSubById(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("subscription id is required")
	}

	err := s.repo.DeleteSubById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) UpdateSubById(ctx context.Context, e entity.Subscription) error {
	if e.Id <= 0 {
		return errors.New("subscription id is required")
	}

	if e.ServiceName == "" {
		return errors.New("service name is required")
	}

	if e.Price < 0 {
		return errors.New("price should be non-negative")
	}

	err := isDateValid(&e.StartDate, e.EndDate)
	if err != nil {
		return err
	}

	err = s.repo.UpdateSubById(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) GetTotalCost(ctx context.Context, userId uuid.UUID, serviceName string, fromDate string, toDate *string) (int, error) {
	if fromDate != "" || toDate != nil {
		err := isDateValid(&fromDate, toDate)
		if err != nil {
			return 0, err
		}
	}

	return s.repo.GetTotalCost(ctx, userId, serviceName, fromDate, toDate)
}

func isDateValid(sourceStartDateStr *string, sourceEndDateStr *string) error {
	startDateStr, startDate, err := parseMMYYYYToFullDate(*sourceStartDateStr)
	if err != nil {
		return err
	}
	*sourceStartDateStr = startDateStr

	if sourceEndDateStr != nil {
		endDateStr, endDate, err := parseMMYYYYToFullDate(*sourceEndDateStr)
		if err != nil {
			return err
		}
		*sourceEndDateStr = endDateStr

		if !startDate.Before(endDate) {
			return errors.New("start date should be before end date")
		}
	}

	return nil
}

func parseMMYYYYToFullDate(dateStr string) (string, time.Time, error) {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 2 {
		return "", time.Time{}, fmt.Errorf("неверный формат: ожидается MM-YYYY")
	}

	month, year := parts[0], parts[1]

	if len(month) != 2 || len(year) != 4 {
		return "", time.Time{}, fmt.Errorf("неверная длина: месяц должен быть 2 цифры, год - 4 цифры")
	}

	fullDateStr := fmt.Sprintf("%s-%s-01", year, month)
	date, err := time.Parse("2006-01-02", fullDateStr)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("неверная дата: %s", dateStr)
	}

	return fullDateStr, date, nil
}
