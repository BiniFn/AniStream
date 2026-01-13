-- Description: This migration removes the unique constraint on the 'mal_id' column in the 'animes' table
--              and drops the foreign key constraint on the 'mal_id' column in the 'anime_metadata' table.
ALTER TABLE anime_metadata
    DROP CONSTRAINT IF EXISTS fk_anime_metadata_mal_id;

ALTER TABLE animes
    DROP CONSTRAINT IF EXISTS unique_mal_id;

CREATE INDEX idx_animes_mal_id ON animes (mal_id);

CREATE INDEX idx_animes_anilist_id ON animes (anilist_id);

