package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hungaikev/rootd/backend/internal/models"
)

// CreateWorkflow handles the creation of a new workflow.
// @Summary Create a new workflow
// @Description A user would typically create a workflow and link it to an existing form schema. The workflow starts in a "draft" status.
// @Tags Workflows
// @Accept  json
// @Produce  json
// @Param   workflow     body    models.Workflow     true        "Workflow to create"
// @Success 201 {object} models.Workflow
// @Router /api/v1/workflows [post]
func CreateWorkflow(c *gin.Context) {
	var newWorkflow models.Workflow
	if err := c.ShouldBindJSON(&newWorkflow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set server-side fields
	newWorkflow.ID = uuid.New().String()
	newWorkflow.Status = models.WorkflowStatusDraft
	newWorkflow.CreatedAt = time.Now()
	newWorkflow.UpdatedAt = time.Now()

	// TODO: Get OwnerID from authenticated user context
	// newWorkflow.OwnerID = ...

	// TODO: Persist the new workflow to the database
	// if err := db.Create(&newWorkflow).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
	// 	return
	// }

	c.JSON(http.StatusCreated, newWorkflow)
}

// ListWorkflows handles listing all workflows for the authenticated user.
// @Summary List all workflows for the authenticated user
// @Description Retrieves a summary list of all workflows owned by the user.
// @Tags Workflows
// @Produce  json
// @Success 200 {array} models.Workflow
// @Router /api/v1/workflows [get]
func ListWorkflows(c *gin.Context) {
	// TODO: Get OwnerID from authenticated user context
	// ownerID := ...

	// TODO: Fetch workflows from the database for the owner
	// var workflows []models.Workflow
	// if err := db.Where("owner_id = ?", ownerID).Find(&workflows).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve workflows"})
	// 	return
	// }

	// Mock response
	workflows := []models.Workflow{
		{ID: uuid.New().String(), Name: "Q3 Customer NPS", Status: models.WorkflowStatusActive},
		{ID: uuid.New().String(), Name: "Job Application", Status: models.WorkflowStatusDraft},
	}

	c.JSON(http.StatusOK, workflows)
}

// GetWorkflow handles retrieving a single workflow.
// @Summary Retrieves a single workflow
// @Description Fetches the complete details of a specific workflow, including its trigger, actions, and submission summary.
// @Tags Workflows
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Success 200 {object} models.Workflow
// @Router /api/v1/workflows/{workflowId} [get]
func GetWorkflow(c *gin.Context) {
	workflowId := c.Param("workflowId")

	// TODO: Fetch workflow from the database
	// var workflow models.Workflow
	// if err := db.First(&workflow, "id = ?", workflowId).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
	// 	return
	// }

	// Mock response
	workflow := models.Workflow{
		ID:        workflowId,
		Name:      "Q3 Customer NPS",
		Status:    models.WorkflowStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, workflow)
}

// UpdateWorkflow handles updating a workflow's configuration.
// @Summary Updates a workflow's configuration
// @Description Used to modify the name, trigger, or actions of a workflow. This operation is typically only allowed when the workflow is in a "draft" status.
// @Tags Workflows
// @Accept  json
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Param   workflow     body    models.Workflow     true        "Updated workflow object"
// @Success 200 {object} models.Workflow
// @Router /api/v1/workflows/{workflowId} [put]
func UpdateWorkflow(c *gin.Context) {
	workflowId := c.Param("workflowId")
	var updatedWorkflow models.Workflow
	if err := c.ShouldBindJSON(&updatedWorkflow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Fetch workflow from the database
	// var workflow models.Workflow
	// if err := db.First(&workflow, "id = ?", workflowId).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
	// 	return
	// }

	// TODO: Check if workflow is in "draft" status
	// if workflow.Status != models.WorkflowStatusDraft {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Workflow can only be updated in draft status"})
	// 	return
	// }

	updatedWorkflow.ID = workflowId
	updatedWorkflow.UpdatedAt = time.Now()

	// TODO: Update workflow in the database
	// if err := db.Save(&updatedWorkflow).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow"})
	// 	return
	// }

	c.JSON(http.StatusOK, updatedWorkflow)
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
func UpdateWorkflowStatus(c *gin.Context) {
	workflowId := c.Param("workflowId")
	var statusUpdate struct {
		Status models.WorkflowStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Fetch workflow from the database
	// var workflow models.Workflow
	// if err := db.First(&workflow, "id = ?", workflowId).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
	// 	return
	// }

	// workflow.Status = statusUpdate.Status
	// workflow.UpdatedAt = time.Now()

	// TODO: Update workflow in the database
	// if err := db.Save(&workflow).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow status"})
	// 	return
	// }

	// Mock response
	workflow := models.Workflow{
		ID:        workflowId,
		Name:      "Q3 Customer NPS",
		Status:    statusUpdate.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
func DeleteWorkflow(c *gin.Context) {
	// workflowId := c.Param("workflowId")

	// TODO: Delete workflow and associated submissions from the database
	// if err := db.Where("id = ?", workflowId).Delete(&models.Workflow{}).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workflow"})
	// 	return
	// }

	c.Status(http.StatusNoContent)
}

// SubmitForm handles the public endpoint for submitting a form.
// @Summary The public endpoint for submitting a form
// @Description This is the URL that users will be sent to. The backend will check if the corresponding workflow's status is "active" before accepting the submission. It is not authenticated.
// @Tags Submissions
// @Accept  json
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Param   submission     body    object     true        "Submission data"
// @Success 201 {object} object
// @Router /w/{workflowId}/submit [post]
func SubmitForm(c *gin.Context) {
	workflowId := c.Param("workflowId")
	var submissionData struct {
		Data     interface{}               `json:"data"`
		Metadata models.SubmissionMetadata `json:"metadata"`
	}
	if err := c.ShouldBindJSON(&submissionData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Fetch workflow from the database
	// var workflow models.Workflow
	// if err := db.First(&workflow, "id = ?", workflowId).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
	// 	return
	// }

	// TODO: Check if workflow is active
	// if workflow.Status != models.WorkflowStatusActive {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "This form is not currently accepting submissions"})
	// 	return
	// }

	// TODO: Persist the new submission to the database
	// if err := db.Create(&newSubmission).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission"})
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{"message": "Submission successful", "workflowId": workflowId})
}

// ListSubmissions handles listing all submissions for a specific workflow.
// @Summary Lists all submissions for a specific workflow
// @Description An authenticated endpoint for the form owner to retrieve all data collected by a workflow. Should support pagination.
// @Tags Submissions
// @Produce  json
// @Param   workflowId     path    string     true        "Workflow ID"
// @Success 200 {array} models.Submission
// @Router /api/v1/workflows/{workflowId}/submissions [get]
func ListSubmissions(c *gin.Context) {
	// workflowId := c.Param("workflowId")

	// TODO: Fetch submissions from the database for the workflow
	// var submissions []models.Submission
	// if err := db.Where("workflow_id = ?", workflowId).Find(&submissions).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
	// 	return
	// }

	// Mock response
	submissions := []models.Submission{
		{ID: uuid.New().String(), CreatedAt: time.Now()},
		{ID: uuid.New().String(), CreatedAt: time.Now().Add(-time.Hour)},
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
func GetSubmission(c *gin.Context) {
	submissionId := c.Param("submissionId")

	// TODO: Fetch submission from the database
	// var submission models.Submission
	// if err := db.First(&submission, "id = ?", submissionId).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
	// 	return
	// }

	// Mock response
	submission := models.Submission{
		ID:        submissionId,
		CreatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, submission)
}
