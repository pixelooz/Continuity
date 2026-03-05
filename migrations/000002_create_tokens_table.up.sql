CREATE TABLE IF NOT EXISTS tokens
(
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    hash       bytea       NOT NULL,
    scope      TEXT        NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_tokens_user_id ON tokens (user_id);
