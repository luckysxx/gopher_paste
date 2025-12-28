-- db/query.sql

-- name: CreatePaste :one
INSERT INTO pastes (
  short_link, content, language, expires_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetPaste :one
SELECT * FROM pastes
WHERE short_link = $1 LIMIT 1;