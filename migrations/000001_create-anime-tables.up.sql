-- Function to generate a NanoID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE
OR REPLACE FUNCTION generate_nanoid(length INT DEFAULT 21) RETURNS TEXT AS $ $ DECLARE alphabet TEXT := '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-';

random_bytes BYTEA;

random_value TEXT := '';

i INT;

BEGIN -- Generate random bytes
random_bytes := gen_random_bytes(length);

-- Generate the NanoID
FOR i IN 1..length LOOP -- Get a value in the range of the alphabet length
random_value := random_value || substr(
  alphabet,
  (get_byte(random_bytes, i - 1) % length) + 1,
  1
);

END LOOP;

RETURN random_value;

END;

$ $ LANGUAGE plpgsql;

-- Create the animes table with a primary key using the NanoID function
CREATE TABLE animes (
  id VARCHAR(21) PRIMARY KEY DEFAULT generate_nanoid(),
  ename TEXT NOT NULL,
  jname TEXT NOT NULL,
  image_url TEXT NOT NULL,
  genre TEXT NOT NULL,
  hi_anime_id TEXT NOT NULL,
  mal_id INT,
  anilist_id INT,
  last_episode INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT unique_hi_anime_id UNIQUE (hi_anime_id),
  CONSTRAINT unique_mal_id UNIQUE (mal_id),
  CONSTRAINT unique_anilist_id UNIQUE (anilist_id)
);

-- Create an ENUM type for media types
CREATE TYPE media_type AS ENUM (
  'tv',
  'movie',
  'ona',
  'ova',
  'special',
  'tv_special',
  'music',
  'cm',
  'pv',
  'unknown'
);

-- Create Rating ENUM type
CREATE TYPE rating AS ENUM (
  'pg_13',
  'r',
  'r+',
  'g',
  'pg',
  'rx',
  'unknown'
);

-- Create an ENUM type for airing status
CREATE TYPE airing_status AS ENUM (
  'finished_airing',
  'currently_airing',
  'not_yet_aired',
  'unknown'
);

-- Create an ENUM type for source
CREATE TYPE source AS ENUM (
  'other',
  'original',
  'manga',
  '4_koma_manga',
  'web_manga',
  'digital_manga',
  'novel',
  'light_novel',
  'visual_novel',
  'game',
  'card_game',
  'book',
  'picture_book',
  'radio',
  'music'
);

-- Create an ENUM type for season
CREATE TYPE season AS ENUM (
  'winter',
  'spring',
  'summer',
  'fall',
  'unknown'
);

-- Create Anime Metadata table
CREATE TABLE anime_metadata (
  mal_id INT PRIMARY KEY,
  description TEXT,
  main_picture_url TEXT,
  media_type media_type NOT NULL DEFAULT 'unknown',
  rating rating NOT NULL DEFAULT 'unknown',
  airing_status airing_status NOT NULL DEFAULT 'unknown',
  avg_episode_duration INT,
  total_episodes INT,
  studio TEXT,
  rank INT,
  mean FLOAT,
  scoringUsers INT,
  popularity INT,
  airing_start_date TEXT,
  airing_end_date TEXT,
  source source NOT NULL DEFAULT 'other',
  trailer_embed_url TEXT,
  season_year INT,
  season season NOT NULL DEFAULT 'unknown',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_anime_metadata_mal_id FOREIGN KEY (mal_id) REFERENCES animes (mal_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create an index on the media_type column for faster lookups
CREATE INDEX idx_anime_metadata_media_type ON anime_metadata (media_type);

-- Create an index on the rating column for faster lookups
CREATE INDEX idx_anime_metadata_rating ON anime_metadata (rating);

-- Create an index on the airing_status column for faster lookups
CREATE INDEX idx_anime_metadata_airing_status ON anime_metadata (airing_status);

-- Create an index on the source column for faster lookups
CREATE INDEX idx_anime_metadata_source ON anime_metadata (source);

-- Create an index on the season column for faster lookups
CREATE INDEX idx_anime_metadata_season ON anime_metadata (season);

-- Search
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX animes_ename_trgm_idx ON animes USING gin(ename gin_trgm_ops);

CREATE INDEX animes_jname_trgm_idx ON animes USING gin(jname gin_trgm_ops);

ALTER TABLE
  animes
ADD
  COLUMN search_vector TSVECTOR;

UPDATE
  animes
SET
  search_vector = to_tsvector('english', ename || ' ' || jname);

CREATE INDEX anime_search_idx ON animes USING gin(search_vector);

CREATE
OR REPLACE FUNCTION animes_search_trigger() RETURNS TRIGGER AS $ $ BEGIN new.search_vector := to_tsvector('english', new.ename || ' ' || new.jname);

RETURN new;

END $ $ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate BEFORE
INSERT
  OR
UPDATE
  ON animes FOR EACH ROW EXECUTE FUNCTION animes_search_trigger();