-- +goose Up
-- +goose StatementBegin
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS assigned_by UUID REFERENCES users(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tickets DROP COLUMN IF EXISTS assigned_by;
-- +goose StatementEnd
