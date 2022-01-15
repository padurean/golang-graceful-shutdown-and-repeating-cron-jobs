package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// JSONErr ...
type JSONErr struct {
	ErrorMessage string
}

func isHTTPMethodValid(r *http.Request, w http.ResponseWriter, method string) bool {
	valid := r.Method == method
	if !valid {
		writeJSONErrResp(w, http.StatusMethodNotAllowed,
			"%s method is not supported on %s resource", r.Method, r.URL.Path)
	}
	return valid
}

func writeJSONErrResp(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	writeJSONResp(w, http.StatusMethodNotAllowed, &JSONErr{
		ErrorMessage: fmt.Sprintf(msg, args...),
	})
}

func writeJSONResp(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("ERROR: failed to write JSON response: %v", err)
	}
}

// HandleGetHello ...
func HandleGetHello(w http.ResponseWriter, r *http.Request) {
	if !isHTTPMethodValid(r, w, http.MethodGet) {
		return
	}
	writeJSONResp(w, http.StatusOK, map[string]string{"message": "Hello World!"})
}

// HandlePostPing ...
func HandlePostPing(w http.ResponseWriter, r *http.Request) {
	if !isHTTPMethodValid(r, w, http.MethodPost) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	payload := make(map[string]string)
	err := decoder.Decode(&payload)
	if err != nil {
		writeJSONErrResp(w, http.StatusBadRequest, "error parsing request body: %v", err)
		return
	}

	writeJSONResp(w, http.StatusOK, &payload)
}
