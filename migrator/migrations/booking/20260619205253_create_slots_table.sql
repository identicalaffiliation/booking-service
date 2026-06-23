-- +goose Up
CREATE TYPE slots_status AS ENUM ('available', 'booked');

CREATE TABLE IF NOT EXISTS slots (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  room_id UUID NOT NULL REFERENCES rooms(id),
  booking_id UUID REFERENCES bookings(id),
  start_at TIMESTAMPTZ NOT NULL,
  end_at TIMESTAMPTZ NOT NULL,
  status slots_status NOT NULL DEFAULT 'available',
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(room_id, start_at)
);
-- +goose Down
DROP TABLE IF EXISTS slots;
DROP TYPE IF EXISTS slots_status;