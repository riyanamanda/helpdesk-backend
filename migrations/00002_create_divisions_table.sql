-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS divisions (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS divisions_name_unique
ON divisions (LOWER(name));

CREATE INDEX IF NOT EXISTS divisions_active_created_at
ON divisions (is_active, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS divisions;
-- +goose StatementEnd
