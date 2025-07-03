package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lerakravts/lite-calendar/internal/model"
	"github.com/lerakravts/lite-calendar/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

type createEventRequest struct {
	Title        string    `json:"title"`
	UserID       string    `json:"user_id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	NotifyBefore *string   `json:"notify_before,omitempty"` // например: "10 minutes"
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Шаг 1: распарсить JSON-запрос
	var req createEventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Шаг 2: собрать Event (мы знаем, что ID и CreatedAt — наша ответственность)
	event := model.Event{
		ID:                 uuid.New().String(),
		Title:              req.Title,
		UserID:             req.UserID,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		NotifyBefore:       req.NotifyBefore,
		NotificationSentAt: nil,
		CreatedAt:          time.Now().UTC(),
	}

	// Шаг 3: передать в сервис
	err = h.service.CreateEvent(ctx, event)
	if err != nil {
		http.Error(w, "failed to create event", http.StatusInternalServerError)
		return
	}

	// Шаг 4: вернуть клиенту ID созданного события
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id": event.ID,
	})
}
func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 1. Читаем параметры from и to из query
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	// 2. Парсим from (если есть), иначе дефолт = сейчас - 1 день
	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			http.Error(w, "invalid 'from' param", http.StatusBadRequest)
			return
		}
	} else {
		from = time.Now().Add(-24 * time.Hour)
	}

	// 3. Парсим to (если есть), иначе дефолт = сейчас + 7 дней
	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			http.Error(w, "invalid 'to' param", http.StatusBadRequest)
			return
		}
	} else {
		to = time.Now().Add(7 * 24 * time.Hour)
	}

	// 4. Защита от "перепутанных дат"
	if to.Before(from) {
		http.Error(w, "'to' must be after 'from'", http.StatusBadRequest)
		return
	}

	// 5. Получаем события через сервис
	events, err := h.service.ListEvents(ctx, from, to)
	if err != nil {
		http.Error(w, "failed to get events", http.StatusInternalServerError)
		return
	}

	// 6. Возвращаем JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/events/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteEvent(ctx, id.String())
	if err != nil {
		http.Error(w, "failed to delete", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
