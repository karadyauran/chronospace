-- name: CreateUserToken :one
INSERT INTO user_tokens (
    user_id,
    refresh_token,
    refresh_token_expires_at
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserTokenByID :one
SELECT * FROM user_tokens
WHERE id = $1;

-- name: GetUserTokenByRefreshToken :one
SELECT * FROM user_tokens
WHERE refresh_token = $1;

-- name: GetUserTokensByUserID :many
SELECT * FROM user_tokens
WHERE user_id = $1;

-- name: UpdateUserToken :one
UPDATE user_tokens
SET 
    refresh_token = $2,
    refresh_token_expires_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUserToken :exec
DELETE FROM user_tokens
WHERE id = $1;

-- name: DeleteUserTokensByUserID :exec
DELETE FROM user_tokens
WHERE user_id = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM user_tokens
WHERE refresh_token_expires_at < CURRENT_TIMESTAMP;

-- name: CountUserTokens :one
SELECT COUNT(*) FROM user_tokens
WHERE user_id = $1;