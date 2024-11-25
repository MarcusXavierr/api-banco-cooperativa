-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1
FOR UPDATE
LIMIT 1;

-- name: GetLastTenTransactions :many
SELECT * FROM transactions
WHERE user_id = $1
ORDER BY id DESC
LIMIT 10;

-- name: RegisterTransaction :exec
INSERT INTO transactions (
    user_id,
    value,
    type,
    description
) VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: IncreaseUserBalance :exec
UPDATE users
SET balance = balance + $1
WHERE id = $2;

-- name: DecreaseUserBalance :exec
UPDATE users
SET balance = balance - $1
WHERE id = $2;

-- name: UpdateUserBalance :exec
UPDATE users
SET balance = $1
WHERE id = $2;

-- name: VerifyCredentials :one
SELECT * FROM users
WHERE email = $1
AND password = $2
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (name, email, password, credit_limit)
  VALUES ($1, $2, $3, $4);

-- name: RetrieveUserFromEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;
