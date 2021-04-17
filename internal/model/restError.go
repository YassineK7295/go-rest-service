package model

// Used to return an error message as the body of a request object
type RestError struct {
	Message string `json:"message"`
}
