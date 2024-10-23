//go:generate mockery --name=AppNameConverter --output=mocks --outpkg=mocks --case=underscore

package middleware

import "github.com/gin-gonic/gin"

// MockAppNameConverter is a mock implementation of AppNameConverter for testing purposes.
type MockAppNameConverter struct{}

// ConvertKeyToAppName mocks the conversion of keys to application names.
func (m *MockAppNameConverter) ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error) {
	return "mock-app", nil // Mock behavior
}
