package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lerakravts/lite-calendar/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

type createEventRequest struct {
	Title string    `json:"title"`
	Time  time.Time `json:"time"`
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events := h.service.ListEvents()

	//устновим заголовок ответа - Говорим клиенту: «Ответ будет в формате JSON»
	w.Header().Add("Content-Type", "application/json")

	//преобразуем events (map[uuid.UUID]Event) в JSON
	err := json.NewEncoder(w).Encode(events)
	if err != nil {
		http.Error(w, "failed to encode events", http.StatusInternalServerError)
		return
	}
}

// можно оставить в CreateEvent проверку метода "на всякий случай",
// но это уже защита от дурака, тк проверяем в switch
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	//1. проверка метода в mux.HandleFunc
	//2. Парсим JSON из тела запроса
	var req createEventRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	//3. Создаем событие через Service
	id, err := h.service.CreateEvent(req.Title, req.Time)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//4. Возвращаем JSON с id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"id": id.String(),
	})

}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	//1.Получаем id из URL-пути: /events/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/events/")

	//2. парсим в uuid.UUID
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid UUID format", http.StatusBadRequest)
		return
	}
	//3. удаляем через сервис
	h.service.DeleteEvent(id)

	//4. Возвращаем 204 No Content — всё прошло успешно, но ничего не возвращаем
	w.WriteHeader(http.StatusNoContent)
}
