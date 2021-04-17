package user

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	RegisterHandlers(r *mux.Router)
}

type router struct {
	controller Controller
}

// Creates a new user router
func NewRouter(db *sql.DB) Router {
	return router{NewController(db)}
}

// Sets up user routes
func (r router) RegisterHandlers(mr *mux.Router) {
	mr.HandleFunc("/users/{userid}", r.controller.Get).Methods(http.MethodGet)
	mr.HandleFunc("/users", r.controller.Create).Methods(http.MethodPost)
	mr.HandleFunc("/users/{userid}", r.controller.Delete).Methods(http.MethodDelete)
	mr.HandleFunc("/users/{userid}", r.controller.Update).Methods(http.MethodPut)
}
