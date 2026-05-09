-- +goose Up
INSERT INTO divisions (name, is_active) VALUES
('IT', TRUE),
('Rekam Medis', TRUE),
('Poli', TRUE),
('Farmasi', TRUE),
('IGD', TRUE);

-- +goose Down
DELETE FROM divisions
WHERE name IN (
  'IT',
  'Rekam Medis',
  'Poli',
  'Farmasi',
  'IGD',
);
