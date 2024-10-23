package middleware

import (
	"GINRouteDemo/middleware/mocks"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func AddUserNameToContext(username string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user-name-key", username)
		c.Next()
	}
}

func TestAuditLogger_MutationByUser(t *testing.T) {

	// Set up log capture using zerolog
	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	// Temporarily replace the global logger with the test logger and restore it after the test
	originalLogger := log.Logger
	log.Logger = logger
	defer func() { log.Logger = originalLogger }()

	// Set Gin to test mode for running in a test environment
	gin.SetMode(gin.TestMode)

	// Use mockery-generated mock
    mockConverter := new(mocks.AppNameConverter)
    mockConverter.On("ConvertKeyToAppName", mock.Anything, "clients", "config").Return("mock-app", nil)

	// Create a Gin engine and attach the middleware
	r := gin.New()
	r.Use(AddUserNameToContext("test-user"))
	r.Use(AuditLogger("appResources", mockConverter)) // Using the AuditLogger middleware
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// Set up the test context using httptest
	w := httptest.NewRecorder()

	// Create a new GET request for the /test route with a custom header
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("gam-key", "test-api-key")

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response status and body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "hello world", w.Body.String())

	// Check the logs to ensure the correct audit message is logged
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Successfully recorded database mutation by a user.")
}

func TestAuditLogger_MutationByApplication(t *testing.T) {
    // Set up log capture using zerolog
    var logBuf bytes.Buffer
    logger := zerolog.New(&logBuf).With().Timestamp().Logger()

    // Temporarily replace the global logger with the test logger and restore it after the test
    originalLogger := log.Logger
    log.Logger = logger
    defer func() { log.Logger = originalLogger }()

    // Set Gin to test mode for running in a test environment
    gin.SetMode(gin.TestMode)

    // Use mockery-generated mock
    mockConverter := new(mocks.AppNameConverter)
    mockConverter.On("ConvertKeyToAppName", mock.Anything, "clients", "config").Return("mock-app", nil)

    // Create a Gin engine and attach the middleware
    r := gin.New()
    r.Use(AuditLogger("appResources", mockConverter)) // Using the mock converter
    r.GET("/test", func(c *gin.Context) {
        c.String(http.StatusOK, "hello world")
    })

    // Set up the test context using httptest
    w := httptest.NewRecorder()

    // Create a new GET request for the /test route with a custom header
    req, _ := http.NewRequest(http.MethodGet, "/test", nil)
    req.Header.Set("gam-key", "test-api-key")

    // Perform the request
    r.ServeHTTP(w, req)

    // Assert the response status and body
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "hello world", w.Body.String())

    // Check the logs to ensure the correct audit message is logged for the application
    logOutput := logBuf.String()
    assert.Contains(t, logOutput, "Successfully recorded database mutation by an application.")
    assert.Contains(t, logOutput, "mock-app")
}

func TestAuditLogger_FailedToRecord(t *testing.T) {

	// Set up log capture using zerolog
	var logBuf bytes.Buffer
	logger := zerolog.New(&logBuf).With().Timestamp().Logger()

	// Temporarily replace the global logger with the test logger and restore it after the test
	originalLogger := log.Logger
	log.Logger = logger
	defer func() { log.Logger = originalLogger }()

	// Set Gin to test mode for running in a test environment
	gin.SetMode(gin.TestMode)

	// Create a mock converter and set up the expectation for ConvertKeyToAppName to return an error
	mockConverter := new(mocks.AppNameConverter)
	mockConverter.On("ConvertKeyToAppName", mock.Anything, "clients", "config").Return("", fmt.Errorf("failed to convert key to application name"))

	// Create a Gin engine and attach the middleware
	r := gin.New()
	r.Use(AuditLogger("appResources", mockConverter)) // Using the mock converter
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	// Set up the test context using httptest
	w := httptest.NewRecorder()

	// Create a new GET request for the /test route with a custom header
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("gam-key", "test-api-key")

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert the response status and body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "hello world", w.Body.String())

	// Check the logs to ensure that the audit log recording failed
	logOutput := logBuf.String()
	assert.Contains(t, logOutput, "Failed to record Audit Logs")
//	assert.Contains(t, logOutput, "failed to convert key to application name")

	// Assert that the mock expectations were met
	mockConverter.AssertExpectations(t)
}