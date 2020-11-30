/* name: GetUsers :many */
SELECT id, created_at, updated_at, name, email
FROM users
WHERE deleted_at is NULL
ORDER BY id;

/* name: GetUser :one */
SELECT id, created_at, updated_at, name, email
FROM users
WHERE id = ? AND deleted_at is NULL;

/* name: GetUserByName :one */
SELECT *
FROM users
WHERE name = ? AND deleted_at is NULL;

/* name: CreateUser :exec */
INSERT INTO users 
    (name, email, password, created_at, updated_at)
VALUES 
    (?, ?, ?, NOW(), NOW());

/* name: UpdateUserEmail :exec */
UPDATE users
SET email = ?
WHERE id = ?;

/* name: DeleteUser :exec */
UPDATE users
SET deleted_at = NOW()
WHERE id = ?