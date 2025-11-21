package chttp

// ErrResponse defines the type that returns in case of the errors.
type ErrResponse struct {
	Description string `json:"error"`
}
