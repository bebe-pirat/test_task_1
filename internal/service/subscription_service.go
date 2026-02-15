package service

import (
	"context"
	"errors"
	"test_task/internal/entity"
	"test_task/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, e entity.Subscription) error {
	if e.ServiceName == "" {
		return errors.New("service name is required")
	}

	if e.Price < 0 {
		return errors.New("price should be non-negative")
	}

	err := s.repo.CreateSubscription(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) GetSubscriptionById(ctx context.Context, id int) (*entity.Subscription, error) {
	if id < 0 {
		return nil, errors.New("id should be non-negative")
	}

	return s.repo.GetSubscriptionById(ctx, id)
}

func (s *SubscriptionService) GetAllSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	return s.repo.GetAllSubsctiptions(ctx)
}

func (s *SubscriptionService) DeleteSubById(ctx context.Context, id int) error {
	if id < 0 {
		return errors.New("id should be non-negative")
	}

	err := s.repo.DeleteSubById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) UpdateSubById(ctx context.Context, e entity.Subscription) error {
	if e.ServiceName == "" {
		return errors.New("service name is required")
	}

	if e.Price < 0 {
		return errors.New("price should be non-negative")
	}

	if e.Id < 0 {
		return errors.New("id should be non-negative")
	}

	err := s.repo.UpdateSubById(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubscriptionService) GetTotalCost(ctx context.Context) (int, error) {
	return s.repo.GetTotalCost(ctx)
}
