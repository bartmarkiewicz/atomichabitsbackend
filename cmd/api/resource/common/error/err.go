package error

import (
	"fmt"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

var (
	DatabaseConnectionFailed = []byte(`{"error":"Database connection failed"}`)
	UpdateFailure            = []byte(`{"error":"Update failed"}`)
	CreateFailure            = []byte(`{"error":"Could not create entity"}`)
	DeleteFailure            = []byte(`{"error":"Could not delete entity"}`)
	JsonEncodeFailure        = []byte(`{"error":"Could not encode entity to JSON"}`)
	JsonDecodeFailure        = []byte(`{"error":"Could not decode entity from JSON"}`)
	InvalidUrlRequest        = []byte(`{"error":"Invalid request url params"}`)
)

func ServerError(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	writeResponse(reps, w)
}

func NotFound(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusNotFound)
	writeResponse(reps, w)
}

func BadRequest(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusBadRequest)
	writeResponse(reps, w)
}

func ValidationErrors(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	writeResponse(reps, w)
}

func writeResponse(reps []byte, w http.ResponseWriter) {
	_, err := w.Write(reps)
	if err != nil {
		fmt.Println("Error writing response: %s", err)
		return
	}
}
