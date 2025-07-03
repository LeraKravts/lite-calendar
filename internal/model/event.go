package model

import "time"

type Event struct {
	ID                 string     `db:"id"`
	Title              string     `db:"title"`
	UserID             string     `db:"user_id"`
	StartTime          time.Time  `db:"start_time"`
	EndTime            time.Time  `db:"end_time"`
	NotifyBefore       *string    `db:"notify_before"`        // INTERVAL → string
	NotificationSentAt *time.Time `db:"notification_sent_at"` // nullable
	CreatedAt          time.Time  `db:"created_at"`
}
