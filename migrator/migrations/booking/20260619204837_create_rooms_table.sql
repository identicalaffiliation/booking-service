-- +goose Up
CREATE TABLE IF NOT EXISTS rooms (
  id UUID PRIMARY KEY NOT NULL,
  name VARCHAR(100) NOT NULL,
  capacity INTEGER NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(name)
);

-- +goose Down
DROP TABLE IF EXISTS rooms;