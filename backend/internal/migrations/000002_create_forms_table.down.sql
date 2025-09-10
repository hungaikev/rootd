-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_forms_updated_at ON forms;
DROP INDEX IF EXISTS idx_forms_created_at;
DROP INDEX IF EXISTS idx_forms_owner_id;
DROP TABLE IF EXISTS forms;
-- +goose StatementEnd
