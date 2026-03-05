CREATE TABLE IF NOT EXISTS notes
(
    id            uuid PRIMARY KEY,
    user_id       uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    collection_id uuid        NOT NULL REFERENCES collections (id) ON DELETE RESTRICT,
    title         TEXT        NOT NULL,
    content       TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT notes_title_not_empty CHECK (TRIM(title) <> ''),
    CONSTRAINT notes_unique_per_collection UNIQUE (user_id, collection_id, title)
);

CREATE INDEX IF NOT EXISTS idx_notes_collection_id ON notes (collection_id);
