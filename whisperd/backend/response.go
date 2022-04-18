package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondError(w http.ResponseWriter, statusCode int, err error) {
	payload := errorResponse{Error: fmt.Sprintf("%v", err)}
	respondJSON(w, statusCode, payload)
}

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(statusCode)
	data, err := json.Marshal(payload)
	if err != nil {
		// not sure a fatal error isn't a bit drastic, but how else could we make ourselves known?
		log.Fatalf("Fatal error encoding response content")
		return
	}
	w.Write(data)
}
