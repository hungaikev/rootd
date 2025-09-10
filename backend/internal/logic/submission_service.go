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

type submissionService struct {
	queries *db.Queries
}

// NewSubmissionService creates a new submission service
func NewSubmissionService(queries *db.Queries) SubmissionService {
	return &submissionService{
		queries: queries,
	}
}

func (s *submissionService) CreateSubmission(ctx context.Context, req CreateSubmissionRequest) (*models.Submission, error) {
	// Validate business rules
	if err := s.validateCreateSubmission(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if workflow is active
	workflowID, err := uuid.Parse(req.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow ID: %w", err)
	}

	workflow, err := s.queries.GetWorkflow(ctx, pgtype.UUID{Bytes: workflowID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow: %w", err)
	}

	if workflow.Status != string(models.WorkflowStatusActive) {
		return nil, fmt.Errorf("workflow is not active and cannot accept submissions")
	}

	// Convert request to database params
	data, _ := json.Marshal(req.Data)
	metadata, _ := json.Marshal(req.Metadata)

	params := db.CreateSubmissionParams{
		WorkflowID: pgtype.UUID{Bytes: workflowID, Valid: true},
		Data:       data,
		Metadata:   metadata,
		Status:     string(models.SubmissionStatusPending),
	}

	if req.SchemaID != nil {
		schemaID := uuid.MustParse(*req.SchemaID)
		params.SchemaID = pgtype.UUID{Bytes: schemaID, Valid: true}
	}

	// Create submission in database
	submission, err := s.queries.CreateSubmission(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to create submission: %w", err)
	}

	// Convert database model to business model
	return s.dbToModel(*submission), nil
}

func (s *submissionService) GetSubmission(ctx context.Context, id string) (*models.Submission, error) {
	submissionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid submission ID: %w", err)
	}

	submission, err := s.queries.GetSubmission(ctx, pgtype.UUID{Bytes: submissionID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	return s.dbToModel(*submission), nil
}

func (s *submissionService) ListSubmissions(ctx context.Context, workflowID string) ([]*models.Submission, error) {
	workflowUUID, err := uuid.Parse(workflowID)
	if err != nil {
		return nil, fmt.Errorf("invalid workflow ID: %w", err)
	}

	submissions, err := s.queries.ListSubmissions(ctx, pgtype.UUID{Bytes: workflowUUID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list submissions: %w", err)
	}

	result := make([]*models.Submission, len(submissions))
	for i, submission := range submissions {
		result[i] = s.dbToModel(*submission)
	}

	return result, nil
}

func (s *submissionService) ListSubmissionsByOwner(ctx context.Context, ownerID string) ([]*models.Submission, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}

	submissions, err := s.queries.ListSubmissionsByOwner(ctx, pgtype.UUID{Bytes: ownerUUID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list submissions by owner: %w", err)
	}

	result := make([]*models.Submission, len(submissions))
	for i, submission := range submissions {
		result[i] = s.dbToModel(*submission)
	}

	return result, nil
}

func (s *submissionService) UpdateSubmissionStatus(ctx context.Context, id string, status models.SubmissionStatus) (*models.Submission, error) {
	submissionID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid submission ID: %w", err)
	}

	// Validate status
	if err := s.validateSubmissionStatus(status); err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	params := db.UpdateSubmissionStatusParams{
		ID:     pgtype.UUID{Bytes: submissionID, Valid: true},
		Status: string(status),
	}

	submission, err := s.queries.UpdateSubmissionStatus(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to update submission status: %w", err)
	}

	return s.dbToModel(*submission), nil
}

func (s *submissionService) DeleteSubmission(ctx context.Context, id string) error {
	submissionID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid submission ID: %w", err)
	}

	return s.queries.DeleteSubmission(ctx, pgtype.UUID{Bytes: submissionID, Valid: true})
}

// Helper methods
func (s *submissionService) validateCreateSubmission(req CreateSubmissionRequest) error {
	if req.WorkflowID == "" {
		return fmt.Errorf("workflow ID is required")
	}
	if req.Data == nil {
		return fmt.Errorf("submission data is required")
	}
	return nil
}

func (s *submissionService) validateSubmissionStatus(status models.SubmissionStatus) error {
	validStatuses := []models.SubmissionStatus{
		models.SubmissionStatusPending,
		models.SubmissionStatusProcessing,
		models.SubmissionStatusCompleted,
		models.SubmissionStatusFailed,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status: %s", status)
}

func (s *submissionService) dbToModel(submission db.Submission) *models.Submission {
	var data map[string]interface{}
	var metadataData map[string]interface{}

	json.Unmarshal(submission.Data, &data)
	json.Unmarshal(submission.Metadata, &metadataData)

	// Convert metadata to proper structure
	metadata := models.SubmissionMetadata{
		IPAddress: metadataData["ipAddress"].(string),
		UserAgent: metadataData["userAgent"].(string),
		Referrer:  metadataData["referrer"].(string),
	}

	return &models.Submission{
		ID:         submission.ID.String(),
		WorkflowID: submission.WorkflowID.String(),
		SchemaID:   submission.SchemaID.String(),
		Data:       data,
		Metadata:   metadata,
		Status:     models.SubmissionStatus(submission.Status),
		CreatedAt:  submission.CreatedAt,
		UpdatedAt:  submission.UpdatedAt,
	}
}
