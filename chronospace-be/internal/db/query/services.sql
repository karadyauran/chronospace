-- name: CreateService :one
INSERT INTO services (
    name,
    description,
    location, 
    price
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetService :one
SELECT * FROM services
WHERE id = $1;

-- name: ListServices :many
SELECT * FROM services
ORDER BY name;

-- name: UpdateService :one
UPDATE services
SET name = $2,
    description = $3,
    location = $4,
    price = $5
WHERE id = $1
RETURNING *;

-- name: DeleteService :exec
DELETE FROM services
WHERE id = $1;