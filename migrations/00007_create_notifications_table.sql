-- +goose Up
CREATE TABLE IF NOT EXISTS notifications(
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL, -- TICKET_ASSIGNED, FEEDBACK_CLOSED, ETC
    reference_type  VARCHAR(50) NOT NULL, -- TICKET, FEEDBACK, ETC
    reference_id    BIGINT NOT NULL, -- ID of entity
    metadata        JSONB NOT NULL DEFAULT '{}',
    is_read         BOOLEAN NOT NULL DEFAULT FALSE,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS notifications_user_id_idx
ON notifications (user_id);

CREATE INDEX IF NOT EXISTS notifications_user_id_is_read_idx
ON notifications (user_id, is_read);

-- +goose Down
DROP TABLE IF EXISTS notifications;
