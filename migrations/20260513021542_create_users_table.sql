-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password VARCHAR(100) NOT NULL,
    google_id VARCHAR(100) UNIQUE,
    avatar_key TEXT,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL,
    division_id INTEGER NOT NULL REFERENCES divisions(id) ON UPDATE RESTRICT ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON UPDATE RESTRICT ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT users_role_check
    CHECK (
        role IN ('EMPLOYEE', 'ADMIN')
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique
ON users (LOWER(email));

CREATE INDEX IF NOT EXISTS users_active_created_at
ON users (google_id, is_active, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
