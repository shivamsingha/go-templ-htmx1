-- name: CountResetTokensByUser :one
SELECT count(*)
FROM reset_tokens
WHERE user_id = $1
    AND expires_at > CURRENT_TIMESTAMP;
-- name: CreateResetToken :exec
INSERT INTO reset_tokens(token, user_id)
VALUES ($1, $2);
-- name: PopResetToken :one
DELETE FROM reset_tokens
WHERE token = $1
RETURNING user_id, expires_at;