package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// 4. Define the API endpoint group.
	// This makes it easy to version your API in the future (e.g., /api/v1, /api/v2).
	api := router.Group("/api/v1")
	{
		// 5. Define the GET handler for the /forms endpoint.
		// When a GET request is made to /api/v1/forms, this function is executed.
		api.GET("/forms", func(c *gin.Context) {
			// 6. Respond with a JSON object and an HTTP 200 OK status.
			// gin.H is a shortcut for map[string]interface{}.
			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": "Forms endpoint acknowledged",
			})
		})
	}

	// 7. Start the HTTP server and listen on port 9000.
	// The server will run indefinitely until it's stopped.
	router.Run(":9000")
}
