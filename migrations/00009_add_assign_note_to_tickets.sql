-- +goose Up
-- +goose StatementBegin
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS assign_note TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tickets DROP COLUMN IF EXISTS assign_note;
-- +goose StatementEnd
