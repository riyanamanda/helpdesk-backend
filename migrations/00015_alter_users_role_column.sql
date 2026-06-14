-- +goose Up
ALTER TABLE users
    DROP CONSTRAINT users_role_check,
    DROP COLUMN role,
    ADD COLUMN role_id BIGINT REFERENCES roles(id) ON UPDATE RESTRICT ON DELETE RESTRICT;

UPDATE users
SET role_id = (SELECT id FROM roles WHERE code = 'EMPLOYEE');

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

-- +goose Down
ALTER TABLE users
    DROP COLUMN role_id,
    ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'EMPLOYEE',
    ADD CONSTRAINT users_role_check CHECK (role IN ('EMPLOYEE', 'ADMIN'));
