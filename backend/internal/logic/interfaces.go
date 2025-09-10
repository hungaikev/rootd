package logic

import (
	"context"

	"github.com/hungaikev/rootd/backend/internal/models"
)

// WorkflowService defines the interface for workflow business logic
type WorkflowService interface {
	CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (*models.Workflow, error)
	GetWorkflow(ctx context.Context, id string) (*models.Workflow, error)
	ListWorkflows(ctx context.Context, ownerID string) ([]*models.Workflow, error)
	UpdateWorkflow(ctx context.Context, id string, req UpdateWorkflowRequest) (*models.Workflow, error)
	UpdateWorkflowStatus(ctx context.Context, id string, status models.WorkflowStatus) (*models.Workflow, error)
	DeleteWorkflow(ctx context.Context, id string) error
}

// FormService defines the interface for form business logic
type FormService interface {
	CreateForm(ctx context.Context, req CreateFormRequest) (*models.Form, error)
	GetForm(ctx context.Context, id string) (*models.Form, error)
	ListForms(ctx context.Context, ownerID string) ([]*models.Form, error)
	UpdateForm(ctx context.Context, id string, req UpdateFormRequest) (*models.Form, error)
	DeleteForm(ctx context.Context, id string) error
}

// SubmissionService defines the interface for submission business logic
type SubmissionService interface {
	CreateSubmission(ctx context.Context, req CreateSubmissionRequest) (*models.Submission, error)
	GetSubmission(ctx context.Context, id string) (*models.Submission, error)
	ListSubmissions(ctx context.Context, workflowID string) ([]*models.Submission, error)
	ListSubmissionsByOwner(ctx context.Context, ownerID string) ([]*models.Submission, error)
	UpdateSubmissionStatus(ctx context.Context, id string, status models.SubmissionStatus) (*models.Submission, error)
	DeleteSubmission(ctx context.Context, id string) error
}

// Request/Response DTOs
type CreateWorkflowRequest struct {
	Name          string                 `json:"name" validate:"required"`
	Description   string                 `json:"description"`
	OwnerID       string                 `json:"owner_id" validate:"required"`
	SchemaID      *string                `json:"schema_id"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Actions       map[string]interface{} `json:"actions"`
}

type UpdateWorkflowRequest struct {
	Name          *string                `json:"name"`
	Description   *string                `json:"description"`
	SchemaID      *string                `json:"schema_id"`
	TriggerConfig map[string]interface{} `json:"trigger_config"`
	Actions       map[string]interface{} `json:"actions"`
}

type CreateFormRequest struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema" validate:"required"`
	OwnerID     string                 `json:"owner_id" validate:"required"`
}

type UpdateFormRequest struct {
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

type CreateSubmissionRequest struct {
	WorkflowID string                     `json:"workflow_id" validate:"required"`
	SchemaID   *string                    `json:"schema_id"`
	Data       map[string]interface{}     `json:"data" validate:"required"`
	Metadata   *models.SubmissionMetadata `json:"metadata"`
}
