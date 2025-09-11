DROP INDEX IF EXISTS animes_year_idx;

DROP INDEX IF EXISTS animes_season_year_idx;

DROP INDEX IF EXISTS animes_genres_arr_gin;

ALTER TABLE animes
  DROP COLUMN IF EXISTS genres_arr;

DROP FUNCTION IF EXISTS normalize_genres(text);

