-- +goose Up
CREATE TABLE IF NOT EXISTS notifications (
  id UUID PRIMARY KEY NOT NULL,
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS notifications;