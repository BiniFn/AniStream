-- 1) Remove full-text trigger and function
DROP TRIGGER IF EXISTS tsvectorupdate ON animes;

DROP FUNCTION IF EXISTS animes_search_trigger();

-- 2) Drop full-text index & column
DROP INDEX IF EXISTS anime_search_idx;

ALTER TABLE animes
  DROP COLUMN IF EXISTS search_vector;

-- 3) Drop trigram extension and its indexes
DROP INDEX IF EXISTS animes_ename_trgm_idx;

DROP INDEX IF EXISTS animes_jname_trgm_idx;

DROP EXTENSION IF EXISTS pg_trgm;

-- 4) Drop anime_metadata and its lookup indexes
DROP TABLE IF EXISTS anime_metadata;

-- 5) Drop ENUM types
DROP TYPE IF EXISTS season;

DROP TYPE IF EXISTS airing_status;

DROP TYPE IF EXISTS rating;

-- 6) Drop indexes on animes (backing the UNIQUEs, if you left them)
DROP INDEX IF EXISTS idx_animes_hi_anime_id;

DROP INDEX IF EXISTS idx_animes_anilist_id;

DROP INDEX IF EXISTS idx_animes_mal_id;

-- 7) Drop the animes table (this also removes its UNIQUE constraints)
DROP TABLE IF EXISTS animes;

-- 8) Drop NanoID generator function
DROP FUNCTION IF EXISTS generate_nanoid(INT);

-- 9) Optionally drop pgcrypto if no longer needed elsewhere
DROP EXTENSION IF EXISTS pgcrypto;

