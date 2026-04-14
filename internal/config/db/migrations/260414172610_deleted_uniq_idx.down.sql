CREATE UNIQUE INDEX IF NOT EXISTS unique_urls_idx ON url_short (url);

DROP INDEX IF EXISTS unique_urls_not_deleted_idx;