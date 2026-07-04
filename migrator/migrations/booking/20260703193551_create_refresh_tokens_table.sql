-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL,
  expires_at BIGINT NOT NULL,
  revoked BOOLEAN NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
