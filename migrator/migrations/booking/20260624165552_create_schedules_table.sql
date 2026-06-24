-- +goose Up
CREATE TABLE IF NOT EXISTS schedules (
  id UUID PRIMARY KEY NOT NULL,
  room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
  work_day DATE NOT NULL,
  start_work_time TIMESTAMPTZ NOT NULL,
  end_work_time TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(room_id, work_day)
);

-- +goose Down
DROP TABLE IF EXISTS schedules;
