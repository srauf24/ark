package model

// SuccessResponse is a generic success response wrapper
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse is a standard error response wrapper
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}

// DeleteResponse is a response wrapper for delete operations
type DeleteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
