package middleware

import "github.com/gin-gonic/gin"

//go:generate mockery --name=AppNameConverter --output=mocks --outpkg=mocks --case=underscore

// AppNameConverter is an interface for converting keys to application names.
type AppNameConverter interface {
	ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error)
}

// RealAppNameConverter is the real implementation of AppNameConverter.
type RealAppNameConverter struct{}

// ConvertKeyToAppName converts keys to the actual application name.
func (r *RealAppNameConverter) ConvertKeyToAppName(ctx *gin.Context, clients string, config string) (string, error) {
	// Implement the actual conversion logic here.
	// For example, fetch from a database or external service.
	return "real-app", nil
}
