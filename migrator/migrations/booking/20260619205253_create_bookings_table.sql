-- +goose Up
CREATE TYPE booking_status AS ENUM ('active', 'cancelled');

CREATE TABLE IF NOT EXISTS bookings (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id),
  slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,

  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(slot_id)
);

-- +goose Down
DROP TABLE IF EXISTS bookings;
DROP TYPE IF EXISTS booking_status;
