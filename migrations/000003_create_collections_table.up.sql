CREATE TABLE IF NOT EXISTS collections
(
    id         uuid PRIMARY KEY,
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    parent_id  uuid REFERENCES collections (id) ON DELETE CASCADE,

    name       TEXT        NOT NULL,
    is_root    BOOLEAN     NOT NULL DEFAULT FALSE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT collections_name_not_empty CHECK (TRIM(name) <> ''),
    CONSTRAINT collection_unique_per_parent UNIQUE (user_id, parent_id, name)
);

CREATE INDEX IF NOT EXISTS idx_collections_parent_id ON collections (parent_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_collections_one_root_per_user ON collections (user_id) WHERE is_root = TRUE;
