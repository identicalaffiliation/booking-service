-- +goose Up
CREATE TYPE users_role AS ENUM ('client', 'admin');

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY NOT NULL,
  nickname VARCHAR(50) NOT NULL,
  password_hash TEXT NOT NULL,
  role users_role NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS users_role;
