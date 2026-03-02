CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users
(
    id            UUID PRIMARY KEY,
    name          VARCHAR(128)        NOT NULL,
    username      VARCHAR(128) UNIQUE NOT NULL,
    email         CITEXT UNIQUE       NOT NULL,
    password_hash BYTEA               NOT NULL,
    activated     BOOL                NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);
