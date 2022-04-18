package backend

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(h *handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/shout/{id}", h.getOne).Methods(http.MethodGet)
	r.HandleFunc("/api/shouts", h.getAll).Methods(http.MethodGet)
	r.HandleFunc("/api/shouts", h.post).Methods(http.MethodPost)
	return r
}
