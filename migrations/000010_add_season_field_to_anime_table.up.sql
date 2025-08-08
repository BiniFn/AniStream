ALTER TABLE animes
  ADD COLUMN season season NOT NULL DEFAULT 'unknown',
  ADD COLUMN season_year int NOT NULL DEFAULT 0;

