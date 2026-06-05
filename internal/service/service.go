package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type Service struct {
	repo   *repository.Repository
	logger *zap.Logger
}

func NewService(repo *repository.Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) CreateSubscription(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	return s.repo.Create(ctx, req)
}

func (s *Service) GetSubscription(ctx context.Context, id int) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, fmt.Errorf("subscription not found")
	}
	return sub, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, id int, req model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	sub, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, fmt.Errorf("subscription not found")
	}
	return sub, nil
}

func (s *Service) DeleteSubscription(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListSubscriptions(ctx context.Context, userID, serviceName string, page, limit int) ([]model.Subscription, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	subscriptions, total, err := s.repo.List(ctx, userID, serviceName, limit, offset)
	return subscriptions, total, page, err
}

func (s *Service) GetTotalCost(ctx context.Context, userID, serviceName, periodFrom, periodTo string) (int, error) {
	if userID == "" {
		return 0, fmt.Errorf("user_id is required")
	}
	return s.repo.GetTotalCost(ctx, userID, serviceName, periodFrom, periodTo)
}
