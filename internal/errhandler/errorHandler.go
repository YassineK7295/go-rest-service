package errhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"

	"github.com/go-sql-driver/mysql"
)

type RestError struct {
	Message string `json:"error"`
}

// Writes message and status code to handler
func WriteMessage(w http.ResponseWriter, msg string, statusCode int) {
	payload, success := getJson(w, RestError{msg})
	if !success {
		return
	}

	log.Println("error: " + msg)

	// every other error
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, payload)
}

// Writes a particular status code to the response, depending on the error
func Write(w http.ResponseWriter, err error) {
	var status int
	var msg string
	me, success := err.(*mysql.MySQLError)
	if success {
		switch me.Number {
		// constraint conflict
		case 1062:
			status = http.StatusBadRequest
			break
		// entity not found
		case 3000:
			status = http.StatusNotFound
		// other mysql errors
		default:
			status = http.StatusInternalServerError
		}
		msg = me.Message
	} else {
		// everything else
		status = http.StatusInternalServerError
		msg = err.Error()
	}

	WriteMessage(w, msg, status)
}

func getJson(w http.ResponseWriter, restError RestError) (string, bool) {
	payload, err := json.Marshal(restError)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("'{\"inner error\":\"%s\"", err.Error()))
		return "", false
	}
	return string(payload), true
}
