-- name: GetUser :one
SELECT id, password, email_verified FROM users
WHERE email = $1 LIMIT 1;