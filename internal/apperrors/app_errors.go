package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error types for different scenarios
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
	ErrValidation        = errors.New("validation failed")
)

// AppError represents a custom application error
type AppError struct {
	Type       error
	Message    string
	StatusCode int
	Details    map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Type.Error()
}

// New creates a new AppError
func New(errType error, message string) *AppError {
	return &AppError{
		Type:       errType,
		Message:    message,
		StatusCode: getStatusCode(errType),
		Details:    make(map[string]interface{}),
	}
}

// Newf creates a new AppError with formatted message
func Newf(errType error, format string, args ...interface{}) *AppError {
	return &AppError{
		Type:       errType,
		Message:    fmt.Sprintf(format, args...),
		StatusCode: getStatusCode(errType),
		Details:    make(map[string]interface{}),
	}
}

// WithDetail adds a detail to the error
func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	e.Details[key] = value
	return e
}

// getStatusCode returns the appropriate HTTP status code for an error type
func getStatusCode(errType error) int {
	switch errType {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists:
		return http.StatusConflict
	case ErrInvalidInput, ErrValidation:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrDatabaseOperation, ErrInternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// IsType checks if an error is of a specific type
func IsType(err error, errType error) bool {
	if appErr, ok := err.(*AppError); ok {
		return errors.Is(appErr.Type, errType)
	}
	return errors.Is(err, errType)
}
