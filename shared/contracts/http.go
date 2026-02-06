package contracts

// APIResponse is the response structure for the API.
type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

// APIError is the error structure for the API.
type APIError struct {
	Code    int64             `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}
