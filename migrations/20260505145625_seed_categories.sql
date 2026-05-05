-- +goose Up
INSERT INTO categories (name, is_active) VALUES
('Network', TRUE),
('Hardware', TRUE),
('Software', TRUE),
('VPN', TRUE),
('Email', TRUE),
('Printer', TRUE),
('Server', TRUE),
('Database', TRUE),
('Security', TRUE),
('Monitoring', TRUE);

-- +goose Down
DELETE FROM categories
WHERE name IN (
  'Network',
  'Hardware',
  'Software',
  'VPN',
  'Email',
  'Printer',
  'Server',
  'Database',
  'Security',
  'Monitoring'
);