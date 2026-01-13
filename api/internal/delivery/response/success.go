package response

type Success[T any] struct {
	Data T   `json:"data"`
	Meta any `json:"meta,omitempty"`
}
