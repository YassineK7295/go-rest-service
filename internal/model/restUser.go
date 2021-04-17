package model

import (
	"net/http"
	"errors"
)
// Used to return a user and the groups it belongs to as the body of a request object
type RestUser struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserId    string    `json:"userid"`
	Groups    *[]string `json:"groups"`
}

// Validates the user object has all of the required fields
// Returns bad request status code otherwise
func (u RestUser) Validate() (error, int) {
	if u.FirstName == "" || u.LastName == "" || u.UserId == "" {
		return errors.New("first_name, last_name, and userid must all be populated"), http.StatusBadRequest
	}
	return nil, 0
}
