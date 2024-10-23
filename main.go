package main

import (
	"GINRouteDemo/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RealAppNameConverter is the real implementation of the AppNameConverter interface
type RealAppNameConverter struct{}

// ConvertKeyToAppName implements the actual conversion logic
func (r *RealAppNameConverter) ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error) {
	// This would contain the real logic for converting a key to an app name
	return "real-app", nil
}

func main() {

	// Create an instance of the real converter
	appNameConverter := &RealAppNameConverter{}

	auditLogger := middleware.AuditLogger("POST", appNameConverter)

	router := gin.Default()

	// Group for Student-related APIs
	studentGroup := router.Group("/students")
	{
		// Create a new student (POST)
		studentGroup.POST("/", auditLogger, func(c *gin.Context) {
			c.String(http.StatusOK, "Student created")
		})

		// Retrieve all students (GET)
		studentGroup.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "List of students")
		})

		// Retrieve a specific student by ID (GET)
		studentGroup.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.String(http.StatusOK, "Retrieve student with ID %s", id)
		})

		// Update a student by ID (PUT)
		studentGroup.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.String(http.StatusOK, "Student with ID %s updated", id)
		})

		// Delete a student by ID (DELETE)
		studentGroup.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.String(http.StatusOK, "Student with ID %s deleted", id)
		})
	}

	// Group for Class-related APIs
	classGroup := router.Group("/class")
	{
		// Create a new class (POST)
		classGroup.POST("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Class created")
		})

		// Retrieve all classes (GET)
		classGroup.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "List of classes")
		})

		// Retrieve a specific class by ID (GET)
		classGroup.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.String(http.StatusOK, "Retrieve class with ID %s", id)
		})

		// Update a class by ID (PUT)
		classGroup.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.String(http.StatusOK, "Class with ID %s updated", id)
		})

		// Delete a class by ID (DELETE)
		classGroup.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.String(http.StatusOK, "Class with ID %s deleted", id)
		})
	}

	// Start the server on port 8080
	router.Run(":8080")
}
