package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/yassinekhaliqui/go-rest-service/internal/group"
	"github.com/yassinekhaliqui/go-rest-service/internal/user"
	"github.com/yassinekhaliqui/go-rest-service/pkg/mw"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

// Set up DB connection and routes
func (a *App) Initialize(config *Config) error {
	var err error

	a.Db, err = sql.Open(config.DB_TYPE, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))
	if err != nil {
		return err
	}

	a.Router = mux.NewRouter()
	a.Router.Use(mw.LogRequest)
	a.Router.Use(mw.AddJsonContentType)

	userRouter := user.NewRouter(a.Db)
	userRouter.RegisterHandlers(a.Router)

	groupRouter := group.NewRouter(a.Db)
	groupRouter.RegisterHandlers(a.Router)
	return nil
}

// Start the server
func (a *App) Run(addr string) error {
	defer a.Db.Close()

	srv := &http.Server{
		Handler:      a.Router,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println(fmt.Sprintf("listening on address %s...", addr))
	return srv.ListenAndServe()
}
