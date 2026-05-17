-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tickets (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,

    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,

    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    priority VARCHAR(20),

    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    closed_by UUID REFERENCES users(id) ON DELETE SET NULL,

    resolution TEXT,

    assigned_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tickets_status_check
    CHECK (
        status IN ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED')
    ),

    CONSTRAINT tickets_priority_check
    CHECK (
        priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')
    )
);

CREATE INDEX IF NOT EXISTS tickets_category_id_idx
ON tickets (category_id);

CREATE INDEX IF NOT EXISTS tickets_status_idx
ON tickets (status);

CREATE INDEX IF NOT EXISTS tickets_priority_idx
ON tickets (priority);

CREATE INDEX IF NOT EXISTS tickets_created_by_idx
ON tickets (created_by);

CREATE INDEX IF NOT EXISTS tickets_assigned_to_idx
ON tickets (assigned_to);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tickets;
-- +goose StatementEnd
