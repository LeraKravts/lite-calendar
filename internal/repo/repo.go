package repo

import (
	"sync"

	"github.com/google/uuid"
	"github.com/lerakravts/lite-calendar/internal/model"
)

// struct будет инициализирована всего 1 экз
type Repository struct {
	mu     sync.RWMutex
	events map[uuid.UUID]model.Event
}

// конструтктор хранилица
func NewRepository() *Repository {
	return &Repository{
		events: make(map[uuid.UUID]model.Event),
	}
}

// пока ошибку не обрабатываем, тк у нас in-memory map — она не падает и не возвращает ошибки
func (r *Repository) Create(event model.Event) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New()
	r.events[id] = event
	return id, nil
}

/*
можно вернуть ошибку, если id не найден

	if _, ok := r.events[id]; !ok {
		return fmt.Errorf("event with ID %s not found", id)
	}
*/
func (r *Repository) Delete(id uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.events, id)

}

func (r *Repository) List() map[uuid.UUID]model.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[uuid.UUID]model.Event)
	for k, v := range r.events {
		result[k] = v
	}
	return result

}
