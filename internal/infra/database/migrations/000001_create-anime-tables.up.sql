-- Function to generate a NanoID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE OR REPLACE FUNCTION generate_nanoid(length int DEFAULT 21)
  RETURNS text
  AS $$
DECLARE
  alphabet text := '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-';
  random_bytes bytea;
  random_value text := '';
  i int;
BEGIN
  -- Generate random bytes
  random_bytes := gen_random_bytes(length);
  -- Generate the NanoID
  FOR i IN 1..length LOOP
    -- Get a value in the range of the alphabet length
    random_value := random_value || substr(alphabet,(get_byte(random_bytes, i - 1) % length) + 1, 1);
  END LOOP;
  RETURN random_value;
END;
$$
LANGUAGE plpgsql;

-- Create the animes table with a primary key using the NanoID function
CREATE TABLE animes(
  id varchar(21) PRIMARY KEY DEFAULT generate_nanoid(),
  ename text NOT NULL,
  jname text NOT NULL,
  image_url text NOT NULL,
  genre text NOT NULL,
  hi_anime_id text NOT NULL,
  mal_id int,
  anilist_id int,
  last_episode int NOT NULL DEFAULT 0,
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW(),
  CONSTRAINT unique_hi_anime_id UNIQUE (hi_anime_id),
  CONSTRAINT unique_mal_id UNIQUE (mal_id)
);

-- Create Rating ENUM type
CREATE TYPE rating AS ENUM(
  'pg_13',
  'r',
  'r+',
  'g',
  'pg',
  'rx',
  'unknown'
);

-- Create an ENUM type for airing status
CREATE TYPE airing_status AS ENUM(
  'finished_airing',
  'currently_airing',
  'not_yet_aired',
  'unknown'
);

-- Create an ENUM type for season
CREATE TYPE season AS ENUM(
  'winter',
  'spring',
  'summer',
  'fall',
  'unknown'
);

-- Create Anime Metadata table
CREATE TABLE anime_metadata(
  mal_id int PRIMARY KEY,
  description text,
  main_picture_url text,
  media_type text,
  rating rating NOT NULL DEFAULT 'unknown',
  airing_status airing_status NOT NULL DEFAULT 'unknown',
  avg_episode_duration int,
  total_episodes int,
  studio text,
  rank int,
  mean float,
  scoringUsers int,
  popularity int,
  airing_start_date text,
  airing_end_date text,
  source text,
  trailer_embed_url text,
  season_year int,
  season season NOT NULL DEFAULT 'unknown',
  created_at timestamp NOT NULL DEFAULT NOW(),
  updated_at timestamp NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_anime_metadata_mal_id FOREIGN KEY (mal_id) REFERENCES animes(mal_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX animes_ename_trgm_idx ON animes USING gin(ename gin_trgm_ops);

CREATE INDEX animes_jname_trgm_idx ON animes USING gin(jname gin_trgm_ops);

ALTER TABLE animes
  ADD COLUMN search_vector tsvector;

UPDATE
  animes
SET
  search_vector = to_tsvector('english', ename || ' ' || jname);

CREATE INDEX anime_search_idx ON animes USING gin(search_vector);

CREATE OR REPLACE FUNCTION animes_search_trigger()
  RETURNS TRIGGER
  AS $$
BEGIN
  NEW.search_vector := to_tsvector('english', NEW.ename || ' ' || NEW.jname);
  RETURN new;
END
$$
LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate
  BEFORE INSERT OR UPDATE ON animes
  FOR EACH ROW
  EXECUTE FUNCTION animes_search_trigger();

