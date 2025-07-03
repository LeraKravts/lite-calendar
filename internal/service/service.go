package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lerakravts/lite-calendar/internal/model"
	"github.com/lerakravts/lite-calendar/internal/repo"
)

// храним *Repository внутри Service
type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateEvent(ctx context.Context, event model.Event) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 1. Проверка: заголовок не пустой
	if strings.TrimSpace(event.Title) == "" {
		return fmt.Errorf("title cannot be empty")
	}

	// 2. Проверка: время старта и окончания
	if event.EndTime.Before(event.StartTime) {
		return fmt.Errorf("end time must be after start time")
	}

	// 3. (опционально) — длина заголовка
	if len(event.Title) > 255 {
		return fmt.Errorf("title too long")
	}
	return s.repo.CreateEvent(ctx, &event)
}

func (s *Service) ListEvents(ctx context.Context, from, to time.Time) ([]model.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return s.repo.ListEvents(ctx, from, to)
}
func (s *Service) DeleteEvent(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return s.repo.DeleteEvent(ctx, id)
}
