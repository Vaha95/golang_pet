ALTER TABLE url_short ADD COLUMN deleted_at TIMESTAMP NULL;

CREATE UNIQUE INDEX unique_urls_not_deleted_idx IF NOT EXISTS
ON url_short (url) 
WHERE (deleted_at IS NULL);

DROP INDEX IF EXISTS unique_urls_idx;