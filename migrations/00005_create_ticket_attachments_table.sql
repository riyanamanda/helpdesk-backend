-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ticket_attachments (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ticket_id INT NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    file_key TEXT NOT NULL,
    attachment_type VARCHAR(20) NOT NULL,
    uploaded_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT ticket_attachments_type_check
    CHECK (
        attachment_type IN ('REPORT', 'RESOLUTION')
    )
);

CREATE INDEX IF NOT EXISTS ticket_attachments_ticket_id_type_created_at
ON ticket_attachments (ticket_id, attachment_type, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ticket_attachments;
-- +goose StatementEnd