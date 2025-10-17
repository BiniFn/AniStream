CREATE TABLE themes(
  id serial PRIMARY KEY,
  name varchar(100) NOT NULL,
  theme_class varchar(100) NOT NULL,
  description text,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO themes(name, theme_class, description)
  VALUES ('Default', '', 'The default Aniways theme'),
('Teal Dream', 'teal', 'A refreshing teal theme'),
('Amber Minimal', 'amber', 'A minimalistic amber theme'),
('Catpuccin', 'catpuccin', 'A cozy catpuccin theme'),
('Cyberpunk', 'cyberpunk', 'A vibrant cyberpunk theme'),
('Ocean Breeze', 'ocean_breeze', 'A cool ocean breeze theme');

ALTER TABLE settings
  ADD COLUMN theme_id integer NOT NULL DEFAULT 1,
  ADD CONSTRAINT fk_theme FOREIGN KEY (theme_id) REFERENCES themes(id) ON DELETE SET DEFAULT
