-- name: CreateWorkflow :one
INSERT INTO workflows (
    name, description, status, owner_id, schema_id, trigger, actions
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetWorkflow :one
SELECT * FROM workflows 
WHERE id = $1;

-- name: ListWorkflows :many
SELECT * FROM workflows 
WHERE owner_id = $1 
ORDER BY created_at DESC;

-- name: UpdateWorkflow :one
UPDATE workflows 
SET 
    name = $2,
    description = $3,
    schema_id = $4,
    trigger = $5,
    actions = $6,
    updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: UpdateWorkflowStatus :one
UPDATE workflows 
SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1 
RETURNING *;

-- name: DeleteWorkflow :exec
DELETE FROM workflows 
WHERE id = $1;

-- name: GetWorkflowSubmissionSummary :one
SELECT 
    COUNT(DISTINCT s.id) as total_submissions,
    COUNT(DISTINCT s.id) as total_visits, -- For now, same as submissions
    COALESCE(AVG(EXTRACT(EPOCH FROM (s.created_at - s.created_at))), 0) as average_time_to_complete,
    MAX(s.created_at) as last_submission_at
FROM submissions s
WHERE s.workflow_id = $1;
