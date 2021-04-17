package group

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

// Creates a new intance of group router
func NewRouter(db *sql.DB) Router {
	return router{NewController(db)}
}

// Registers the group endpoints with the router
func (r router) RegisterHandlers(mr *mux.Router) {
	mr.HandleFunc("/groups/{groupName}", r.controller.Get).Methods(http.MethodGet)
	mr.HandleFunc("/groups", r.controller.Create).Methods(http.MethodPost)
	mr.HandleFunc("/groups/{groupName}", r.controller.Delete).Methods(http.MethodDelete)
	mr.HandleFunc("/groups/{groupName}", r.controller.Update).Methods(http.MethodPut)
}
