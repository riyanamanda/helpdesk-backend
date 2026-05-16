-- +goose Up
CREATE TABLE IF NOT EXISTS ticket_resolutions (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ticket_id INTEGER NOT NULL UNIQUE REFERENCES tickets(id) ON DELETE CASCADE,
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolution TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS ticket_resolutions_resolved_by_idx
ON ticket_resolutions (resolved_by);

-- +goose Down
DROP TABLE IF EXISTS ticket_resolutions;
