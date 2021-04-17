package model

import (
	"errors"
	"net/http"
)

// Used to return a group name as the body of a request object
type RestGroup struct {
	Name string `json:"name"`
}

// Validates the object has the name field populated
// Returns a bad request status code otherwise
func (r RestGroup) Validate() (error, int) {
	if r.Name == "" {
		return errors.New("name field must be populated"), http.StatusBadRequest
	}
	return nil, 0
}

// Used to return a list of users in a group as the body of a request object
type RestGroupMembers struct {
	UserIds *[]string `json:"userids"`
}
