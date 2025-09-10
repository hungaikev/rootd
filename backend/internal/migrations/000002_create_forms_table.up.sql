-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS forms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    schema JSONB NOT NULL,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_forms_owner_id ON forms(owner_id);
CREATE INDEX IF NOT EXISTS idx_forms_created_at ON forms(created_at);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_forms_updated_at 
    BEFORE UPDATE ON forms 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd
