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
-- name: CountResetTokensByUser :one
SELECT count(*)
FROM reset_tokens
WHERE user_id = $1
    AND expires_at > CURRENT_TIMESTAMP;
-- name: CreateResetToken :exec
INSERT INTO reset_tokens(token, user_id)
VALUES ($1, $2);