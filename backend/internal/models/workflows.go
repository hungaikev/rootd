package models

import (
	"encoding/json"
	"time"
)

// WorkflowStatus represents the lifecycle state of a form's workflow.
type WorkflowStatus string

const (
	WorkflowStatusDraft    WorkflowStatus = "draft"
	WorkflowStatusActive   WorkflowStatus = "active"
	WorkflowStatusPaused   WorkflowStatus = "paused"
	WorkflowStatusStopped  WorkflowStatus = "stopped"
	WorkflowStatusArchived WorkflowStatus = "archived"
)

// TriggerType defines what initiates a workflow.
type TriggerType string

const (
	TriggerTypeFormSubmission TriggerType = "form_submission"
	TriggerTypeWebhook        TriggerType = "webhook"
	TriggerTypeManual         TriggerType = "manual"
)

// ActionType defines the type of operation a workflow node performs.
type ActionType string

const (
	ActionTypeSendEmail    ActionType = "send_email"
	ActionTypeCallWebhook  ActionType = "call_webhook"
	ActionTypeNotification ActionType = "notification"
	ActionTypeCondition    ActionType = "condition"
)

// Workflow represents the operational controller for a form schema.
// It defines the trigger, manages the state, and contains the sequence of actions to be executed.
type Workflow struct {
	ID        string         `json:"id"`        // UUID for the workflow.
	Name      string         `json:"name"`      // User-defined name for the workflow (e.g., "Q3 Customer NPS").
	OwnerID   string         `json:"ownerId"`   // The user who owns this workflow.
	SchemaID  string         `json:"schemaId"`  // The ID of the form schema this workflow controls.
	Status    WorkflowStatus `json:"status"`    // The current state of the workflow.
	Trigger   Trigger        `json:"trigger"`   // The event that starts this workflow.
	Actions   []Action       `json:"actions"`   // The sequence of steps to execute.
	CreatedAt time.Time      `json:"createdAt"` // Timestamp of creation.
	UpdatedAt time.Time      `json:"updatedAt"` // Timestamp of last update.

	// SubmissionSummary holds aggregated data about the submissions for this workflow.
	SubmissionSummary SubmissionSummary `json:"submissionSummary"`
}

// Trigger defines the event that initiates a workflow.
type Trigger struct {
	Type   TriggerType     `json:"type"`             // The type of trigger.
	Config json.RawMessage `json:"config,omitempty"` // Configuration specific to the trigger type (e.g., webhook URL).
}

// Action represents a single step or node within a workflow.
type Action struct {
	ID          string          `json:"id"`                    // UUID for this action step.
	Type        ActionType      `json:"type"`                  // The type of action to perform.
	Description string          `json:"description,omitempty"` // A user-defined description of the step.
	Config      json.RawMessage `json:"config"`                // Configuration for the action (e.g., email template, webhook URL, conditions).
	Conditional *Conditional    `json:"conditional,omitempty"` // Optional logic to determine if this action should run.
}

// SubmissionSummary contains aggregated analytics for a workflow.
type SubmissionSummary struct {
	TotalVisits           int        `json:"totalVisits"`                     // Total number of times the form was viewed.
	TotalSubmissions      int        `json:"totalSubmissions"`                // Total number of submissions received.
	CompletionRate        float64    `json:"completionRate"`                  // Percentage of visits that resulted in a submission.
	AverageTimeToComplete int        `json:"averageTimeToComplete,omitempty"` // Average time in seconds from visit to submission.
	LastSubmissionAt      *time.Time `json:"lastSubmissionAt,omitempty"`      // Timestamp of the most recent submission.
}

// Submission represents a single data entry for a form through an active workflow.
type Submission struct {
	ID         string             `json:"id"`         // UUID for the submission.
	WorkflowID string             `json:"workflowId"` // The ID of the workflow it belongs to.
	SchemaID   string             `json:"schemaId"`   // The ID of the schema used for this submission.
	CreatedAt  time.Time          `json:"createdAt"`  // Timestamp of submission.
	Data       json.RawMessage    `json:"data"`       // The submitted form data, stored as raw JSON.
	Metadata   SubmissionMetadata `json:"metadata"`   // Additional metadata about the submission context.
}

// SubmissionMetadata contains contextual information about a submission.
type SubmissionMetadata struct {
	IPAddress string `json:"ipAddress,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	Referrer  string `json:"referrer,omitempty"`
}
