CREATE TABLE IF NOT EXISTS url_short (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url VARCHAR(3000) NOT NULL,
    short VARCHAR(20) UNIQUE NOT NULL,
    ext_id VARCHAR(255) NOT NULL
);

COMMENT ON COLUMN url_short.ext_id IS 'Id to url from external system.';