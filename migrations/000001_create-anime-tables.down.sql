-- Drop trigger and its function
DROP TRIGGER IF EXISTS tsvectorupdate ON animes;
DROP FUNCTION IF EXISTS animes_search_trigger();
-- Drop GIN search indexes
DROP INDEX IF EXISTS anime_search_idx;
DROP INDEX IF EXISTS animes_ename_trgm_idx;
DROP INDEX IF EXISTS animes_jname_trgm_idx;
-- Remove the tsvector column
ALTER TABLE animes DROP COLUMN IF EXISTS search_vector;
-- Drop trigram extension (optional—only if not used elsewhere)
DROP EXTENSION IF EXISTS pg_trgm;
-- Remove foreign key relationship
ALTER TABLE anime_metadata DROP CONSTRAINT IF EXISTS fk_anime_metadata_mal_id;
-- Drop indexes on anime_metadata
DROP INDEX IF EXISTS idx_anime_metadata_season;
DROP INDEX IF EXISTS idx_anime_metadata_source;
DROP INDEX IF EXISTS idx_anime_metadata_airing_status;
DROP INDEX IF EXISTS idx_anime_metadata_rating;
DROP INDEX IF EXISTS idx_anime_metadata_media_type;
-- Drop tables (automatically drops their indexes, but we’ve cleared FKs first)
DROP TABLE IF EXISTS anime_metadata;
DROP TABLE IF EXISTS animes;
-- Drop NanoID generator
DROP FUNCTION IF EXISTS generate_nanoid(INT);
-- Drop ENUM types
DROP TYPE IF EXISTS season;
DROP TYPE IF EXISTS source;
DROP TYPE IF EXISTS airing_status;
DROP TYPE IF EXISTS rating;
DROP TYPE IF EXISTS media_type;