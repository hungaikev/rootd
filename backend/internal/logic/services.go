package logic

import (
	"github.com/hungaikev/rootd/backend/internal/db"
)

// Services holds all the business logic services
type Services struct {
	Workflow   WorkflowService
	Form       FormService
	Submission SubmissionService
}

// NewServices creates a new services container
func NewServices(queries *db.Queries) *Services {
	return &Services{
		Workflow:   NewWorkflowService(queries),
		Form:       NewFormService(queries),
		Submission: NewSubmissionService(queries),
	}
}
