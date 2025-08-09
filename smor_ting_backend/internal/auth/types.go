package auth

// ErrorResponse represents a standard error payload for auth endpoints in this package
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
