-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetLastTenTransactions :many
SELECT * FROM transactions
WHERE user_id = $1 LIMIT 10;
