ALTER TABLE animes
    ADD CONSTRAINT unique_mal_id UNIQUE (mal_id);

ALTER TABLE anime_metadata
    ADD CONSTRAINT fk_anime_metadata_mal_id FOREIGN KEY (mal_id) REFERENCES animes (mal_id) ON DELETE CASCADE;

DROP INDEX IF EXISTS idx_animes_mal_id;

DROP INDEX IF EXISTS idx_animes_anilist_id;

