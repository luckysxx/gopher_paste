BEGIN;

ALTER TABLE pastes
  ADD COLUMN IF NOT EXISTS id BIGINT,
  ADD COLUMN IF NOT EXISTS owner_id BIGINT,
  ADD COLUMN IF NOT EXISTS title VARCHAR(120),
  ADD COLUMN IF NOT EXISTS visibility VARCHAR(20),
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP;

CREATE SEQUENCE IF NOT EXISTS pastes_id_seq;
ALTER TABLE pastes ALTER COLUMN id SET DEFAULT nextval('pastes_id_seq');

UPDATE pastes
SET id = nextval('pastes_id_seq')
WHERE id IS NULL;

SELECT setval('pastes_id_seq', COALESCE((SELECT MAX(id) FROM pastes), 1), true);

UPDATE pastes
SET owner_id = 0
WHERE owner_id IS NULL;

UPDATE pastes
SET title = COALESCE(NULLIF(short_link, ''), 'Untitled Snippet')
WHERE title IS NULL;

UPDATE pastes
SET visibility = 'private'
WHERE visibility IS NULL;

UPDATE pastes
SET updated_at = created_at
WHERE updated_at IS NULL;

ALTER TABLE pastes ALTER COLUMN id SET NOT NULL;
ALTER TABLE pastes ALTER COLUMN owner_id SET NOT NULL;
ALTER TABLE pastes ALTER COLUMN title SET NOT NULL;
ALTER TABLE pastes ALTER COLUMN visibility SET NOT NULL;
ALTER TABLE pastes ALTER COLUMN visibility SET DEFAULT 'private';
ALTER TABLE pastes ALTER COLUMN updated_at SET NOT NULL;
ALTER TABLE pastes ALTER COLUMN updated_at SET DEFAULT NOW();

ALTER TABLE pastes DROP CONSTRAINT IF EXISTS pastes_pkey;
ALTER TABLE pastes ADD CONSTRAINT pastes_pkey PRIMARY KEY (id);

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'pastes_short_link_key'
  ) THEN
    ALTER TABLE pastes ADD CONSTRAINT pastes_short_link_key UNIQUE (short_link);
  END IF;
END $$;

ALTER TABLE pastes DROP COLUMN IF EXISTS expires_at;

CREATE INDEX IF NOT EXISTS idx_pastes_owner_id ON pastes(owner_id);

COMMIT;
