-- name: CreateForm :one
INSERT INTO forms (
    name, description, schema, owner_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetForm :one
SELECT * FROM forms 
WHERE id = $1;

-- name: ListForms :many
SELECT * FROM forms 
WHERE owner_id = $1 
ORDER BY created_at DESC;

-- name: UpdateForm :one
UPDATE forms 
SET 
    name = $2,
    description = $3,
    schema = $4,
    updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: DeleteForm :exec
DELETE FROM forms 
WHERE id = $1;
