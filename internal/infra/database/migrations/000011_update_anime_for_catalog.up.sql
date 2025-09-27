CREATE OR REPLACE FUNCTION normalize_genres(csv text)
  RETURNS text[]
  LANGUAGE sql
  IMMUTABLE
  RETURNS NULL ON NULL INPUT
  AS $$
  SELECT
    COALESCE((
      SELECT
        ARRAY( SELECT DISTINCT
            lower(trim(x))
        FROM unnest(string_to_array(csv, ',')) AS t(x)
        WHERE
          trim(x) <> '' ORDER BY lower(trim(x)))), ARRAY[]::text[]);
$$;

ALTER TABLE animes
  ADD COLUMN genres_arr text[] GENERATED ALWAYS AS (normalize_genres(genre)) STORED;

CREATE INDEX animes_genres_arr_gin ON animes USING GIN(genres_arr);

CREATE INDEX animes_season_year_idx ON animes(season, season_year);

CREATE INDEX animes_year_idx ON animes(season_year);

