CREATE UNIQUE INDEX unique_urls_not_deleted_idx IF NOT EXISTS
ON url_short (url) 
WHERE (deleted_at IS NULL);

CREATE INDEX IF NOT EXISTS short_by_user_idx
ON url_short (short, created_by_user)
WHERE (deleted_at IS NULL);

DROP INDEX IF EXISTS unique_urls_idx;