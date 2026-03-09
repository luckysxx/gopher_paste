-- db/schema.sql

CREATE TABLE IF NOT EXISTS pastes (
    id          BIGSERIAL PRIMARY KEY,
    owner_id    BIGINT NOT NULL,
    title       VARCHAR(120) NOT NULL,
    short_link  VARCHAR(10) UNIQUE,
    content     TEXT NOT NULL,
    language    VARCHAR(20) NOT NULL DEFAULT 'text',
    visibility  VARCHAR(20) NOT NULL DEFAULT 'private',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pastes_owner_id ON pastes(owner_id);