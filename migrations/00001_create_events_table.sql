-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE events (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        title TEXT NOT NULL,
                        user_id UUID NOT NULL,
                        start_time TIMESTAMP NOT NULL,
                        end_time TIMESTAMP NOT NULL,
                        notify_before INTERVAL,
                        notification_sent_at TIMESTAMP,
                        created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd

