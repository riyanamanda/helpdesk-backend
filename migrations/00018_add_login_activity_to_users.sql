-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_method VARCHAR(10),
    ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMPTZ;

ALTER TABLE users
    ADD CONSTRAINT users_last_login_method_check
    CHECK (
        last_login_method IN ('PASSWORD', 'GOOGLE')
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_last_login_method_check,
    DROP COLUMN IF EXISTS last_active_at,
    DROP COLUMN IF EXISTS last_login_method,
    DROP COLUMN IF EXISTS last_login_at;
-- +goose StatementEnd
