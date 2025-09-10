-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_workflows_updated_at ON workflows;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_workflows_created_at;
DROP INDEX IF EXISTS idx_workflows_status;
DROP INDEX IF EXISTS idx_workflows_owner_id;
DROP TABLE IF EXISTS workflows;
-- +goose StatementEnd
