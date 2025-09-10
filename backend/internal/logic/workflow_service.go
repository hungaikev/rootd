package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/hungaikev/rootd/backend/internal/db"
	"github.com/hungaikev/rootd/backend/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type workflowService struct {
	queries *db.Queries
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(queries *db.Queries) WorkflowService {
	return &workflowService{
		queries: queries,
	}
}

func (s *workflowService) CreateWorkflow(ctx context.Context, req CreateWorkflowRequest) (*models.Workflow, error) {
	// Validate business rules
	if err := s.validateCreateWorkflow(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert request to database params
	trigger, _ := json.Marshal(map[string]interface{}{
		"type":   "manual", // Default trigger type
		"config": req.TriggerConfig,
	})
	actions, _ := json.Marshal(req.Actions)

	ownerUUID := uuid.MustParse(req.OwnerID)
	params := db.CreateWorkflowParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Status:      string(models.WorkflowStatusDraft),
		OwnerID:     pgtype.UUID{Bytes: ownerUUID, Valid: true},
		Trigger:     trigger,
		Actions:     actions,
	}

	if req.SchemaID != nil {
		schemaID := uuid.MustParse(*req.SchemaID)
		params.SchemaID = pgtype.UUID{Bytes: schemaID, Valid: true}
	}

	// Create workflow in database
	workflow, err := s.queries.CreateWorkflow(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	// Convert database model to business model
	return s.dbToModel(*workflow), nil
}

func (s *workflowService) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
	workflowID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow ID: %w", err)
	}

	workflow, err := s.queries.GetWorkflow(ctx, pgtype.UUID{Bytes: workflowID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	return s.dbToModel(*workflow), nil
}

func (s *workflowService) ListWorkflows(ctx context.Context, ownerID string) ([]*models.Workflow, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}

	workflows, err := s.queries.ListWorkflows(ctx, pgtype.UUID{Bytes: ownerUUID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list workflows: %w", err)
	}

	result := make([]*models.Workflow, len(workflows))
	for i, workflow := range workflows {
		result[i] = s.dbToModel(*workflow)
	}

	return result, nil
}

func (s *workflowService) UpdateWorkflow(ctx context.Context, id string, req UpdateWorkflowRequest) (*models.Workflow, error) {
	workflowID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow ID: %w", err)
	}

	// Get existing workflow to check if it's in draft status
	existing, err := s.queries.GetWorkflow(ctx, pgtype.UUID{Bytes: workflowID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get existing workflow: %w", err)
	}

	// Business rule: Only allow updates if workflow is in draft status
	if existing.Status != string(models.WorkflowStatusDraft) {
		return nil, fmt.Errorf("workflow can only be updated in draft status")
	}

	// Prepare update params
	params := db.UpdateWorkflowParams{
		ID: pgtype.UUID{Bytes: workflowID, Valid: true},
	}

	if req.Name != nil {
		params.Name = *req.Name
	} else {
		params.Name = existing.Name
	}

	if req.Description != nil {
		params.Description = pgtype.Text{String: *req.Description, Valid: true}
	} else {
		params.Description = existing.Description
	}

	if req.SchemaID != nil {
		schemaID := uuid.MustParse(*req.SchemaID)
		params.SchemaID = pgtype.UUID{Bytes: schemaID, Valid: true}
	} else {
		params.SchemaID = existing.SchemaID
	}

	if req.TriggerConfig != nil {
		trigger, _ := json.Marshal(map[string]interface{}{
			"type":   "manual", // Default trigger type
			"config": req.TriggerConfig,
		})
		params.Trigger = trigger
	} else {
		params.Trigger = existing.Trigger
	}

	if req.Actions != nil {
		actions, _ := json.Marshal(req.Actions)
		params.Actions = actions
	} else {
		params.Actions = existing.Actions
	}

	// Update workflow in database
	workflow, err := s.queries.UpdateWorkflow(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to update workflow: %w", err)
	}

	return s.dbToModel(*workflow), nil
}

func (s *workflowService) UpdateWorkflowStatus(ctx context.Context, id string, status models.WorkflowStatus) (*models.Workflow, error) {
	workflowID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow ID: %w", err)
	}

	// Validate status transition
	if err := s.validateStatusTransition(status); err != nil {
		return nil, fmt.Errorf("invalid status transition: %w", err)
	}

	params := db.UpdateWorkflowStatusParams{
		ID:     pgtype.UUID{Bytes: workflowID, Valid: true},
		Status: string(status),
	}

	workflow, err := s.queries.UpdateWorkflowStatus(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to update workflow status: %w", err)
	}

	return s.dbToModel(*workflow), nil
}

func (s *workflowService) DeleteWorkflow(ctx context.Context, id string) error {
	workflowID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid workflow ID: %w", err)
	}

	// Check if workflow has submissions
	submissions, err := s.queries.ListSubmissions(ctx, pgtype.UUID{Bytes: workflowID, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to check workflow submissions: %w", err)
	}

	if len(submissions) > 0 {
		return fmt.Errorf("cannot delete workflow with existing submissions")
	}

	return s.queries.DeleteWorkflow(ctx, pgtype.UUID{Bytes: workflowID, Valid: true})
}

// Helper methods
func (s *workflowService) validateCreateWorkflow(req CreateWorkflowRequest) error {
	if req.Name == "" {
		return fmt.Errorf("workflow name is required")
	}
	if req.OwnerID == "" {
		return fmt.Errorf("owner ID is required")
	}
	return nil
}

func (s *workflowService) validateStatusTransition(status models.WorkflowStatus) error {
	validStatuses := []models.WorkflowStatus{
		models.WorkflowStatusDraft,
		models.WorkflowStatusActive,
		models.WorkflowStatusPaused,
		models.WorkflowStatusArchived,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status: %s", status)
}

func (s *workflowService) dbToModel(workflow db.Workflow) *models.Workflow {
	var triggerData map[string]interface{}
	var actions []map[string]interface{}

	json.Unmarshal(workflow.Trigger, &triggerData)
	json.Unmarshal(workflow.Actions, &actions)

	// Convert actions to the proper structure
	modelActions := make([]models.Action, len(actions))
	for i, action := range actions {
		config, _ := json.Marshal(action["config"])
		modelActions[i] = models.Action{
			ID:          action["id"].(string),
			Type:        models.ActionType(action["type"].(string)),
			Description: action["description"].(string),
			Config:      config,
		}
	}

	// Convert trigger
	triggerConfig, _ := json.Marshal(triggerData["config"])
	trigger := models.Trigger{
		Type:   models.TriggerType(triggerData["type"].(string)),
		Config: triggerConfig,
	}

	schemaID := ""
	if workflow.SchemaID.Valid {
		schemaID = uuid.UUID(workflow.SchemaID.Bytes[:]).String()
	}

	return &models.Workflow{
		ID:        uuid.UUID(workflow.ID.Bytes[:]).String(),
		Name:      workflow.Name,
		Status:    models.WorkflowStatus(workflow.Status),
		OwnerID:   uuid.UUID(workflow.OwnerID.Bytes[:]).String(),
		SchemaID:  schemaID,
		Trigger:   trigger,
		Actions:   modelActions,
		CreatedAt: workflow.CreatedAt,
		UpdatedAt: workflow.UpdatedAt,
	}
}
