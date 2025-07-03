package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lerakravts/lite-calendar/internal/model"
)

// struct будет инициализирована всего 1 экз
type Repository struct {
	db *sqlx.DB
}

// конструтктор хранилица
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateEvent(ctx context.Context, event *model.Event) error {
	query := `
		INSERT INTO events (
			id, title, user_id, start_time, end_time,
			notify_before, notification_sent_at, created_at
		) VALUES (
			:id, :title, :user_id, :start_time, :end_time,
			:notify_before, :notification_sent_at, :created_at
		);
	`

	_, err := r.db.NamedExecContext(ctx, query, event)
	return err
}

func (r *Repository) ListEvents(ctx context.Context, from, to time.Time) ([]model.Event, error) {
	query := `SELECT id, title, user_id, start_time, end_time,
                     notify_before, notification_sent_at, created_at
              FROM events
              WHERE start_time < $1 AND end_time > $2
              ORDER BY start_time`

	var events []model.Event
	err := r.db.SelectContext(ctx, &events, query, to, from)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	return events, nil
}

func (r *Repository) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}

	//если события не было
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("event not found")
	}
	return nil
}
