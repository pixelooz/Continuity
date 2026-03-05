CREATE TABLE IF NOT EXISTS collections
(
    id         uuid PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    parent_id  uuid REFERENCES collections (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT collections_name_not_empty CHECK (TRIM(name) <> ''),
    CONSTRAINT collection_unique_per_parent UNIQUE (user_id, parent_id, name)
);

CREATE INDEX IF NOT EXISTS idx_collections_parent_id ON collections (parent_id);
