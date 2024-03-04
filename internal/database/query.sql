-- name: GetUser :one
SELECT id,
    password,
    email_verified
FROM users
WHERE email = $1
LIMIT 1;
-- name: CreateUser :exec
INSERT INTO users(email, name, password)
VALUES ($1, $2, $3);