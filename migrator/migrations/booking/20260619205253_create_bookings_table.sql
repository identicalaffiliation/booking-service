-- +goose Up
CREATE TYPE booking_status AS ENUM ('active', 'cancelled');

CREATE TABLE IF NOT EXISTS bookings (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
  status booking_status NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX
    unique_slot_id_by_active_booking
ON
    bookings(slot_id)
WHERE
    status = 'active';

-- +goose Down
DROP INDEX IF EXISTS unique_slot_id_by_active_booking;
DROP TABLE IF EXISTS bookings;
DROP TYPE IF EXISTS booking_status;
