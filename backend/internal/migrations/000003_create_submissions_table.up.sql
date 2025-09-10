-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    schema_id UUID REFERENCES forms(id) ON DELETE SET NULL,
    data JSONB NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_submissions_workflow_id ON submissions(workflow_id);
CREATE INDEX IF NOT EXISTS idx_submissions_schema_id ON submissions(schema_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status);
CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_submissions_updated_at 
    BEFORE UPDATE ON submissions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd
