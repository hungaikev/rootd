-- name: CreateSubmission :one
INSERT INTO submissions (
    workflow_id, schema_id, data, metadata, status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSubmission :one
SELECT * FROM submissions 
WHERE id = $1;

-- name: ListSubmissions :many
SELECT * FROM submissions 
WHERE workflow_id = $1 
ORDER BY created_at DESC;

-- name: ListSubmissionsByOwner :many
SELECT s.* FROM submissions s
JOIN workflows w ON s.workflow_id = w.id
WHERE w.owner_id = $1 
ORDER BY s.created_at DESC;

-- name: UpdateSubmissionStatus :one
UPDATE submissions 
SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: DeleteSubmission :exec
DELETE FROM submissions 
WHERE id = $1;
