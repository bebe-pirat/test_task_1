package service

import (
	"context"
	"errors"
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

func (s *SubscriptionService) CreateSubscription(ctx context.Context, e entity.Subscription) (uuid.UUID, error) {
	if e.ServiceName == "" {
		return uuid.Nil, errors.New("service name is required")
	}

	if e.Price < 0 {
		return uuid.Nil, errors.New("price should be non-negative")
	}

	err := isDateValid(e.StartDate, e.EndDate)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := s.repo.CreateSubscription(ctx, e)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *SubscriptionService) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*entity.Subscription, error) {
	if id == uuid.Nil {
		return nil, errors.New("subscription id is required")
	}

	return s.repo.GetSubscriptionById(ctx, id)
}

func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	return s.repo.GetAllSubscriptions(ctx)
}

func (s *SubscriptionService) DeleteSubById(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("subscription id is required")
	}

	err := s.repo.DeleteSubById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) UpdateSubById(ctx context.Context, e entity.Subscription) error {
	if e.Id == uuid.Nil {
		return errors.New("subscription id is required")
	}

	if e.ServiceName == "" {
		return errors.New("service name is required")
	}

	if e.Price < 0 {
		return errors.New("price should be non-negative")
	}

	err := isDateValid(e.StartDate, e.EndDate)
	if err != nil {
		return err
	}

	err = s.repo.UpdateSubById(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) GetTotalCost(ctx context.Context, subId uuid.UUID, serviceName string, fromDate string, toDate *string) (int, error) {
	if fromDate != "" && toDate != nil {
		err := isDateValid(fromDate, toDate)
		if err != nil {
			return 0, err
		}
	}

	return s.repo.GetTotalCost(ctx, subId, serviceName, fromDate, toDate)
}

func isDateValid(startDateStr string, endDateStr *string) error {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return err
	}
	if endDateStr != nil {
		endDate, err := time.Parse("2006-01-02", *endDateStr)
		if err != nil {
			return err
		}

		if !startDate.Before(endDate) {
			return errors.New("start date should be before end date")
		}
	}
	return nil
}
