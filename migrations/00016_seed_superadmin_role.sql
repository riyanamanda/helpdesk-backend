-- +goose Up
INSERT INTO roles (code) VALUES ('SUPERADMIN');

-- +goose Down
DELETE FROM roles WHERE code = 'SUPERADMIN';
