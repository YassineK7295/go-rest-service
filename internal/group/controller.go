package group

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yassinekhaliqui/go-rest-service/internal/errhandler"
	"github.com/yassinekhaliqui/go-rest-service/internal/model"
	"github.com/yassinekhaliqui/go-rest-service/pkg/util"
)

type Controller interface {
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	service Service
}

// Creates new controller instance
func NewController(db *sql.DB) Controller {
	return controller{NewService(db)}
}

// Retrieves a list of users that are part of the group
// Returns 404 if group is not found
func (a controller) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["groupName"]

	group, users, err := a.service.GetWithUsers(r.Context(), groupName)
	if err != nil {
		errhandler.Write(w, err)
		return
	}

	if group == (model.Group{}) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("group %s not found\n", groupName)))
		return
	}

	restGroupMembers := toRestGroupMembers(users)
	respBody, err := json.Marshal(restGroupMembers)
	if err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(respBody))
}

// Creates an empty group
// Returns 400 if group already exists
func (a controller) Create(w http.ResponseWriter, r *http.Request) {
	var restGroup model.RestGroup
	if err := json.NewDecoder(r.Body).Decode(&restGroup); err != nil {
		errhandler.Write(w, err)
		return
	}
	defer r.Body.Close()

	if err, statusCode := restGroup.Validate(); err != nil {
		errhandler.WriteMessage(w, err.Error(), statusCode)
		return
	}

	_, err := a.service.Insert(r.Context(), model.Group{Name: restGroup.Name})
	if err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("group %s has been created\n", restGroup.Name)))
}

// Deletes a group and any links to users for that group
// Returns 404 if group is not found
func (a controller) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["groupName"]

	if err := a.service.Delete(r.Context(), groupName); err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("group %s has been deleted\n", groupName)))
}

// Updates group membership
// Returns 404 if group is not found
func (a controller) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["groupName"]

	var restGroupMembers model.RestGroupMembers
	if err := json.NewDecoder(r.Body).Decode(&restGroupMembers); err != nil {
		errhandler.Write(w, err)
		return
	}
	defer r.Body.Close()

	if err := a.service.UpdateGroupMembership(r.Context(), groupName, restGroupMembers.UserIds); err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("group %s has been updated\n", groupName)))
}

// Converts a Group object to a RestGroup object
func toRestGroup(group model.Group) model.RestGroup {
	return model.RestGroup{group.Name}
}

// Converts an array of users to a RestGroupMembers object
func toRestGroupMembers(users *[]model.User) model.RestGroupMembers {
	var userIds []string

	for _, user := range *users {
		userIds = append(userIds, user.UserId)
	}

	return model.RestGroupMembers{&userIds}
}
