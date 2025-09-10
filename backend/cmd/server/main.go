package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hungaikev/rootd/backend/internal/api/handlers"
	"github.com/hungaikev/rootd/backend/internal/db"
	"github.com/hungaikev/rootd/backend/internal/logic"
)

func main() {
	// Get database configuration from environment variables
	cfg := db.ServiceConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "rootd"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Create database service
	dbService, err := db.NewService(cfg)
	if err != nil {
		log.Fatal("Failed to create database service:", err)
	}
	defer dbService.Close()

	// Run migrations
	if err := dbService.RunMigrations(nil); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Create business logic services
	services := logic.NewServices(dbService.Queries)

	// Create handlers
	workflowHandlers := handlers.NewWorkflowHandlers(services)

	// Initialize Gin router with default middleware (logger, recovery)
	router := gin.Default()

	// Configure CORS middleware to allow frontend communication
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:8787"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	apiV1 := router.Group("/api/v1")
	{
		// Workflow Management Endpoints
		workflows := apiV1.Group("/workflows")
		{
			workflows.POST("", workflowHandlers.CreateWorkflow)
			workflows.GET("", workflowHandlers.ListWorkflows)
			workflows.GET("/:workflowId", workflowHandlers.GetWorkflow)
			workflows.PUT("/:workflowId", workflowHandlers.UpdateWorkflow)
			workflows.PATCH("/:workflowId/status", workflowHandlers.UpdateWorkflowStatus)
			workflows.DELETE("/:workflowId", workflowHandlers.DeleteWorkflow)
			workflows.GET("/:workflowId/submissions", workflowHandlers.ListSubmissions)
		}

		// Submission Management Endpoints
		submissions := apiV1.Group("/submissions")
		{
			submissions.GET("/:submissionId", workflowHandlers.GetSubmission)
		}
	}

	// Public submission endpoint
	public := router.Group("/w")
	{
		public.POST("/:workflowId/submit", workflowHandlers.SubmitForm)
	}

	// Start the HTTP server
	port := getEnv("PORT", "9000")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
