CREATE TABLE IF NOT EXISTS tags
(
    id         uuid PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tags_name_not_empty CHECK (TRIM(name) <> ''),
    CONSTRAINT tags_unique_per_user UNIQUE (user_id, name)
);
