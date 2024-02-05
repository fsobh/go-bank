-- name: CreateTransfer :one
INSERT INTO transfers(
  from_account_id, to_account_id, amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransferByID :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetAllTransfersByFromAccountID :many
SELECT * FROM transfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetAllTransfersByToAccountID :many
SELECT * FROM transfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetAllTransfersByBetween :many
SELECT * FROM transfers
WHERE to_account_id = $1
AND   from_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;