package user

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

// Creates a new instance of the user controller
func NewController(db *sql.DB) Controller {
	return controller{
		service: NewService(db),
	}
}

// Gets a user and the groups they belong to
// Returns 404 if user is not found
func (a controller) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]

	user, groups, err := a.service.GetWithGroup(r.Context(), userId)
	if err != nil {
		errhandler.Write(w, err)
		return
	}

	if user == (model.User{}) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("user id %s was not found", userId)))
		return
	}

	restUser := merge(user, groups)
	payload, err := json.Marshal(restUser)
	if err != nil {
		errhandler.Write(w, err)
		return
	}

	fmt.Fprintf(w, string(payload))
}

// Creates a new user with any groups (if provided)
// Returns 400 if userid is duplicated
func (a controller) Create(w http.ResponseWriter, r *http.Request) {
	var restUser model.RestUser
	if err := json.NewDecoder(r.Body).Decode(&restUser); err != nil {
		errhandler.Write(w, err)
		return
	}
	defer r.Body.Close()

	if err, statusCode := restUser.Validate(); err != nil {
		errhandler.WriteMessage(w, err.Error(), statusCode)
		return
	}

	user, groupNames := deconstruct(restUser)

	if err := a.service.InsertTx(r.Context(), user, groupNames); err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("user %s created\n", restUser.UserId)))
}

// Deletes a user and their linkages to groups
// Returns 404 if not found
func (a controller) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]

	if err := a.service.Delete(r.Context(), userId); err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("user %s has been deleted\n", userId)))
}

// Updates a user and their linkages to groups
// Returns 404 if user is not found
func (a controller) Update(w http.ResponseWriter, r *http.Request) {
	var restUser model.RestUser
	if err := json.NewDecoder(r.Body).Decode(&restUser); err != nil {
		errhandler.Write(w, err)
		return
	}
	defer r.Body.Close()

	if err, statusCode := restUser.Validate(); err != nil {
		errhandler.WriteMessage(w, err.Error(), statusCode)
		return
	}

	user, groupNames := deconstruct(restUser)

	if err := a.service.UpdateTx(r.Context(), user, groupNames); err != nil {
		errhandler.Write(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, util.MessageJson("result", fmt.Sprintf("user %s has been updated\n", restUser.UserId)))
}

// Creates a RestUser from a User and an array of Groups
func merge(user model.User, groups *[]model.Group) model.RestUser {
	groupNames := make([]string, len(*groups))

	for i, g := range *groups {
		groupNames[i] = g.Name
	}

	return model.RestUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserId:    user.UserId,
		Groups:    &groupNames,
	}
}

// Breaks the RestUser object up into a User obj and an array of group names
func deconstruct(restUser model.RestUser) (model.User, *[]string) {
	return model.User{
		FirstName: restUser.FirstName,
		LastName:  restUser.LastName,
		UserId:    restUser.UserId,
	}, restUser.Groups
}
