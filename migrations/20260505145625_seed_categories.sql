-- +goose Up
INSERT INTO categories (name, is_active) VALUES
('Network', TRUE),
('Hardware', TRUE),
('Software', TRUE);

-- +goose Down
DELETE FROM categories
WHERE name IN (
  'Network',
  'Hardware',
  'Software',
);