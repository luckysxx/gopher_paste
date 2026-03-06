-- db/schema.sql

CREATE TABLE IF NOT EXISTS pastes (
    short_link  VARCHAR(10) PRIMARY KEY,
    content     TEXT NOT NULL,
    language    VARCHAR(20) NOT NULL DEFAULT 'text',
    expires_at  TIMESTAMP,  --如果不填就是永久
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);