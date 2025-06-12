package service

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lerakravts/lite-calendar/internal/model"
	"github.com/lerakravts/lite-calendar/internal/repo"
)

// храним *Repository внутри Service
type Service struct {
	repo *repo.Repository
}

func NewService(r *repo.Repository) *Service {
	return &Service{repo: r}
}

// пока ctx не прокидываем
// Здесь CreateEvent ничего не добавляет — просто проксирует результат repo.Create.
func (s *Service) CreateEvent(title string, eventTime time.Time) (uuid.UUID, error) {
	if strings.TrimSpace(title) == "" {
		return uuid.Nil, errors.New("title can't be empty")
	}
	if eventTime.Before(time.Now()) {
		return uuid.Nil, errors.New("event time can't be in the past")
	}

	event := model.Event{
		Title: title,
		Time:  eventTime,
	}

	return s.repo.Create(event)
}

func (s *Service) DeleteEvent(id uuid.UUID) {
	s.repo.Delete(id)
}

func (s *Service) ListEvents() map[uuid.UUID]model.Event {
	return s.repo.List()
}
