package errors

import "fmt"

type AppError struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

// Add Error method to implement the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Status, e.Message)
}

func NewAppError(statusCode int, status string, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
	}
}
