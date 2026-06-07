-- +goose Up
CREATE TABLE IF NOT EXISTS feedbacks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,

    type VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',

    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,

    reviewed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT feedbacks_type_check
    CHECK (
        type IN ('FEATURE_REQUEST', 'IMPROVEMENT', 'BUG_REPORT')
    ),

    CONSTRAINT feedbacks_status_check
    CHECK (
        status IN ('OPEN', 'IN_REVIEW', 'ACCEPTED', 'REJECTED', 'DELIVERED')
    )
);

CREATE INDEX IF NOT EXISTS feedbacks_type_idx
ON feedbacks (type);

CREATE INDEX IF NOT EXISTS feedbacks_status_idx
ON feedbacks (status);

CREATE INDEX IF NOT EXISTS feedbacks_created_by_idx
ON feedbacks (created_by);

CREATE INDEX IF NOT EXISTS feedbacks_reviewed_by_idx
ON feedbacks (reviewed_by);

-- +goose Down
DROP TABLE IF EXISTS feedbacks;