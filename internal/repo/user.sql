-- name: GetUser :one
SELECT * FROM users WHERE guid = @guid;

-- name: InsertUser :exec
INSERT INTO users (guid, name, occupation, email, pwd, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, now(), now());

-- name: UpdateUser :exec
UPDATE users SET name = @name, occupation = @occupation, updated_at = now() WHERE guid = @guid;

-- name: DeleteUser :exec
UPDATE users SET is_deleted = true, updated_at = now() WHERE guid = @guid;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = @email;

