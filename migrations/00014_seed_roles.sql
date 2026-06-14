-- +goose Up
INSERT INTO roles (code) VALUES ('ADMIN'), ('EMPLOYEE');

-- +goose Down
DELETE FROM roles WHERE code IN ('ADMIN', 'EMPLOYEE');
