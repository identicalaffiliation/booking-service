-- +goose Up
CREATE TABLE IF NOT EXISTS slots (
  id UUID PRIMARY KEY NOT NULL,
  room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  start_time INT NOT NULL CHECK (start_time >= 0 AND start_time < 1440),
  end_time INT NOT NULL CHECK (end_time > start_time AND end_time < 1440),
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(room_id, date, start_time)
);
-- +goose Down
DROP TABLE IF EXISTS slots;