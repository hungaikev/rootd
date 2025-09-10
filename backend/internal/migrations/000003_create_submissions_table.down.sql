-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_submissions_updated_at ON submissions;
DROP INDEX IF EXISTS idx_submissions_created_at;
DROP INDEX IF EXISTS idx_submissions_status;
DROP INDEX IF EXISTS idx_submissions_schema_id;
DROP INDEX IF EXISTS idx_submissions_workflow_id;
DROP TABLE IF EXISTS submissions;
-- +goose StatementEnd
