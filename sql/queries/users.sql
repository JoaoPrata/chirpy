-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users u
INNER JOIN refresh_tokens t ON u.id = t.user_id
WHERE t.token = $1
LIMIT 1;

-- name: UpdateUserCredentials :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpdateUserChirpyRed :one
UPDATE users
SET is_chirpy_red = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;