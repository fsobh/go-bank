-- name: CreateEntry :one
INSERT INTO entries(
  account_id, amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEntryByID :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntry :many
SELECT * FROM entries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetAllEntriesByAccountID :many
SELECT * FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;