package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hungaikev/rootd/backend/internal/api/handlers"
)

func main() {
	// 1. Initialize a new Gin router with default middleware (logger, recovery).
	router := gin.Default()

	// 2. Configure CORS middleware to allow frontend communication
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API v1 group
	apiV1 := router.Group("/api/v1")
	{
		// Workflow Management Endpoints
		workflows := apiV1.Group("/workflows")
		{
			workflows.POST("", handlers.CreateWorkflow)
			workflows.GET("", handlers.ListWorkflows)
			workflows.GET("/:workflowId", handlers.GetWorkflow)
			workflows.PUT("/:workflowId", handlers.UpdateWorkflow)
			workflows.PATCH("/:workflowId/status", handlers.UpdateWorkflowStatus)
			workflows.DELETE("/:workflowId", handlers.DeleteWorkflow)
			workflows.GET("/:workflowId/submissions", handlers.ListSubmissions)
		}

		// Submission Management Endpoints
		submissions := apiV1.Group("/submissions")
		{
			submissions.GET("/:submissionId", handlers.GetSubmission)
		}
	}

	// Public submission endpoint
	public := router.Group("/w")
	{
		public.POST("/:workflowId/submit", handlers.SubmitForm)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 7. Start the HTTP server and listen on port 9000.
	// The server will run indefinitely until it's stopped.
	router.Run(":9000")
}
