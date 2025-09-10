package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hungaikev/rootd/backend/internal/logic"
	"github.com/hungaikev/rootd/backend/internal/models"
)

// UpdateWorkflow handles updating a workflow's configuration.
// @Summary Updates a workflow's configuration
// @Description Used to modify the name, trigger, or actions of a workflow. This operation is typically only allowed when the workflow is in a "draft" status.
// @Tags Workflows
// @Accept  json
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Param   workflow     body    logic.UpdateWorkflowRequest     true        "Updated workflow object"
// @Success 200 {object} models.Workflow
// @Router /api/v1/workflows/{workflowId} [put]
func (h *WorkflowHandlers) UpdateWorkflow(c *gin.Context) {
	workflowID := c.Param("workflowId")
	var req logic.UpdateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.services.Workflow.UpdateWorkflow(c.Request.Context(), workflowID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// UpdateWorkflowStatus handles changing the status of a workflow.
// @Summary Changes the status of a workflow
// @Description A dedicated endpoint to manage the workflow's state (e.g., from "draft" to "active", or "active" to "paused").
// @Tags Workflows
// @Accept  json
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Param   status     body    object     true        "New status"
// @Success 200 {object} models.Workflow
// @Router /api/v1/workflows/{workflowId}/status [patch]
func (h *WorkflowHandlers) UpdateWorkflowStatus(c *gin.Context) {
	workflowID := c.Param("workflowId")
	var statusUpdate struct {
		Status models.WorkflowStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflow, err := h.services.Workflow.UpdateWorkflowStatus(c.Request.Context(), workflowID, statusUpdate.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// DeleteWorkflow handles deleting a workflow.
// @Summary Deletes a workflow
// @Description Permanently deletes a workflow and all of its associated submissions. This is a destructive action.
// @Tags Workflows
// @Param   workflowId     path    string     true        "Workflow ID"
// @Success 204 {object} nil
// @Router /api/v1/workflows/{workflowId} [delete]
func (h *WorkflowHandlers) DeleteWorkflow(c *gin.Context) {
	workflowID := c.Param("workflowId")

	err := h.services.Workflow.DeleteWorkflow(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SubmitForm handles the public endpoint for submitting a form.
// @Summary The public endpoint for submitting a form
// @Description This is the URL that users will be sent to. The backend will check if the corresponding workflow's status is "active" before accepting the submission. It is not authenticated.
// @Tags Submissions
// @Accept  json
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Param   submission     body    logic.CreateSubmissionRequest     true        "Submission data"
// @Success 201 {object} object
// @Router /w/{workflowId}/submit [post]
func (h *WorkflowHandlers) SubmitForm(c *gin.Context) {
	workflowID := c.Param("workflowId")
	var req logic.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.WorkflowID = workflowID

	submission, err := h.services.Submission.CreateSubmission(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Submission successful",
		"workflowId":   workflowID,
		"submissionId": submission.ID,
	})
}

// ListSubmissions handles listing all submissions for a specific workflow.
// @Summary Lists all submissions for a specific workflow
// @Description An authenticated endpoint for the form owner to retrieve all data collected by a workflow. Should support pagination.
// @Tags Submissions
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Success 200 {array} models.Submission
// @Router /api/v1/workflows/{workflowId}/submissions [get]
func (h *WorkflowHandlers) ListSubmissions(c *gin.Context) {
	workflowID := c.Param("workflowId")

	submissions, err := h.services.Submission.ListSubmissions(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// GetSubmission handles retrieving a single submission.
// @Summary Retrieves a single submission
// @Description An authenticated endpoint to get the full details of one specific submission, including its data and metadata.
// @Tags Submissions
// @Produce  json
// @Param   submissionId     path    string     true        "Submission ID"
// @Success 200 {object} models.Submission
// @Router /api/v1/submissions/{submissionId} [get]
func (h *WorkflowHandlers) GetSubmission(c *gin.Context) {
	submissionID := c.Param("submissionId")

	submission, err := h.services.Submission.GetSubmission(c.Request.Context(), submissionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	c.JSON(http.StatusOK, submission)
}
