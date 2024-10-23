package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	zLog "github.com/rs/zerolog/log"
)

// AuditLogger middleware captures and logs HTTP request and response details, including
func AuditLogger(appResources string, converter AppNameConverter) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Record the request path
		httpRequest := fmt.Sprintf("%s: %s", ctx.Request.Method, ctx.Request.URL.Path)

		// Proceed with the request
		ctx.Next()

		// Record the response status after processing the request
		responseStatus := strconv.Itoa(ctx.Writer.Status())

		// Log all request headers
		for key, values := range ctx.Request.Header {
			fmt.Printf("%s: %s\n", key, values)
		}

		// Log context keys set using ctx.Set()
		for _, key := range ctx.Keys {
			fmt.Printf("%v: %v\n", key, ctx.GetString(key.(string)))
		}

		// Prepare the base logger with common fields
		logEntry := zLog.Info().
			Str("tag", "audit").
			Str("source", httpRequest).
			Str("response_status", responseStatus)

		// Check for userName in context and log if found
		if userName, userExists := ctx.Get("user-name-key"); userExists {
			logEntry.Str("actor", userName.(string)).
				Msg("Successfully recorded database mutation by a user.")
			return
		}

		// If userName does not exist, check for applicationName
		applicationName, appErr := converter.ConvertKeyToAppName(ctx, "clients", "config")
		if appErr != nil {
			// Log a warning if application name cannot be retrieved
			zLog.Warn().
				Str("tag", "audit").
				Str("source", httpRequest).
				Msg("Failed to record Audit Logs.")
			return
		}

		// Log the applicationName if no userName is found
		logEntry.Str("actor", applicationName).
			Msg("Successfully recorded database mutation by an application.")
	}
}

//go:generate mockery --name=AppNameConverter --output=mocks --outpkg=mocks --case=underscore

// AppNameConverter is an interface for converting keys to application names
type AppNameConverter interface {
	ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error)
}

type RealAppNameConverter struct{}

func (r *RealAppNameConverter) ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error) {
	// Actual implementation
	return "test-app", nil
}

type MockAppNameConverter struct{}

func (m *MockAppNameConverter) ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error) {
	return "mock-app", nil // Mock behavior
}
