package response

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
	Meta  any   `json:"meta,omitempty"`
}

func NewError(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Code:    code,
			Message: message,
		},
	}
}

func NewErrorWithDetails(code, message string, details any) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

func NewInternalServerError(message string) ErrorResponse {
	return NewError("internal_server_error", message)
}

var (
	ErrInvalidRequest = NewError("invalid_request", "The request is invalid.")
	ErrUnauthorized   = NewError("unauthorized", "unauthorized")
)
