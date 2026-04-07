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
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM users u
INNER JOIN refresh_tokens t ON u.id = t.user_id
WHERE t.token = $1;

-- name: UpdateUser :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;