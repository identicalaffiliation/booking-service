-- +goose Up
CREATE TABLE IF NOT EXISTS slots (
  id UUID PRIMARY KEY NOT NULL,
  room_id UUID NOT NULL REFERENCES rooms(id),
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(room_id, start_time, end_time)
);
-- +goose Down
DROP TABLE IF EXISTS slots;