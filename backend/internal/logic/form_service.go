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

type formService struct {
	queries *db.Queries
}

// NewFormService creates a new form service
func NewFormService(queries *db.Queries) FormService {
	return &formService{
		queries: queries,
	}
}

func (s *formService) CreateForm(ctx context.Context, req CreateFormRequest) (*models.Form, error) {
	// Validate business rules
	if err := s.validateCreateForm(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Convert request to database params
	schema, _ := json.Marshal(req.Schema)

	ownerUUID := uuid.MustParse(req.OwnerID)
	params := db.CreateFormParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
		Schema:      schema,
		OwnerID:     pgtype.UUID{Bytes: ownerUUID, Valid: true},
	}

	// Create form in database
	form, err := s.queries.CreateForm(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to create form: %w", err)
	}

	// Convert database model to business model
	return s.dbToModel(*form), nil
}

func (s *formService) GetForm(ctx context.Context, id string) (*models.Form, error) {
	formID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid form ID: %w", err)
	}

	form, err := s.queries.GetForm(ctx, pgtype.UUID{Bytes: formID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get form: %w", err)
	}

	return s.dbToModel(*form), nil
}

func (s *formService) ListForms(ctx context.Context, ownerID string) ([]*models.Form, error) {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, fmt.Errorf("invalid owner ID: %w", err)
	}

	forms, err := s.queries.ListForms(ctx, pgtype.UUID{Bytes: ownerUUID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list forms: %w", err)
	}

	result := make([]*models.Form, len(forms))
	for i, form := range forms {
		result[i] = s.dbToModel(*form)
	}

	return result, nil
}

func (s *formService) UpdateForm(ctx context.Context, id string, req UpdateFormRequest) (*models.Form, error) {
	formID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid form ID: %w", err)
	}

	// Get existing form
	existing, err := s.queries.GetForm(ctx, pgtype.UUID{Bytes: formID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get existing form: %w", err)
	}

	// Prepare update params
	params := db.UpdateFormParams{
		ID: pgtype.UUID{Bytes: formID, Valid: true},
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

	if req.Schema != nil {
		schema, _ := json.Marshal(req.Schema)
		params.Schema = schema
	} else {
		params.Schema = existing.Schema
	}

	// Update form in database
	form, err := s.queries.UpdateForm(ctx, &params)
	if err != nil {
		return nil, fmt.Errorf("failed to update form: %w", err)
	}

	return s.dbToModel(*form), nil
}

func (s *formService) DeleteForm(ctx context.Context, id string) error {
	formID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid form ID: %w", err)
	}

	// Check if form is being used by any workflows
	// This would require a query to check workflow references
	// For now, we'll allow deletion

	return s.queries.DeleteForm(ctx, pgtype.UUID{Bytes: formID, Valid: true})
}

// Helper methods
func (s *formService) validateCreateForm(req CreateFormRequest) error {
	if req.Name == "" {
		return fmt.Errorf("form name is required")
	}
	if req.OwnerID == "" {
		return fmt.Errorf("owner ID is required")
	}
	if req.Schema == nil {
		return fmt.Errorf("form schema is required")
	}
	return nil
}

func (s *formService) dbToModel(form db.Form) *models.Form {
	var schema map[string]interface{}
	json.Unmarshal(form.Schema, &schema)

	description := ""
	if form.Description.Valid {
		description = form.Description.String
	}

	return &models.Form{
		ID:          uuid.UUID(form.ID.Bytes[:]).String(),
		Name:        form.Name,
		Description: description,
		Schema:      schema,
		OwnerID:     uuid.UUID(form.OwnerID.Bytes[:]).String(),
		CreatedAt:   form.CreatedAt,
		UpdatedAt:   form.UpdatedAt,
	}
}
