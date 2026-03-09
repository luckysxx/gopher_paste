-- db/query.sql

-- name: CreatePaste :one
INSERT INTO pastes (
  owner_id, title, short_link, content, language, visibility
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetPaste :one
SELECT * FROM pastes
WHERE id = $1 LIMIT 1;

-- name: ListMyPastes :many
SELECT * FROM pastes
WHERE owner_id = $1
ORDER BY updated_at DESC;

-- name: UpdatePaste :one
UPDATE pastes
SET
  title = $2,
  content = $3,
  language = $4,
  visibility = $5,
  updated_at = NOW()
WHERE id = $1 AND owner_id = $6
RETURNING *;