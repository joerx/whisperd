package backend

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"whisperd.io/whisperd/whisperd"
	"whisperd.io/whisperd/whisperd/db"
	"whisperd.io/whisperd/whisperd/db/store"
)

type shoutResponse struct {
	Shout whisperd.Shout `json:"shout"`
}

type shoutListResponse struct {
	Shouts []whisperd.Shout `json:"shouts"`
}

type putShoutRequest struct {
	Shout whisperd.Shout `json:"shout"`
}

type handler struct {
	shouts store.Shouts
}

func newHandler(shouts store.Shouts) *handler {
	return &handler{shouts: shouts}
}

func (b *handler) getAll(w http.ResponseWriter, r *http.Request) {
	sl, err := b.shouts.GetAll(r.Context())

	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
	} else {
		slr := shoutListResponse{sl}
		respondJSON(w, http.StatusOK, slr)
	}
}

func (b *handler) getOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ii, err := strconv.Atoi(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	shout, err := b.shouts.Get(r.Context(), int64(ii))
	switch err {
	case nil:
		respondJSON(w, http.StatusOK, shoutResponse{shout})
	case db.ErrorNotFound:
		respondError(w, http.StatusNotFound, err)
	default:
		respondError(w, http.StatusInternalServerError, err)
	}
}

func (b *handler) post(w http.ResponseWriter, r *http.Request) {
	var payload putShoutRequest

	defer r.Body.Close()
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(bs, &payload); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	sr, err := b.shouts.Insert(r.Context(), payload.Shout)
	switch err {
	case nil:
		respondJSON(w, http.StatusOK, shoutResponse{sr})
	case db.ErrorInvalidRecord:
		respondError(w, http.StatusBadRequest, err)
	default:
		respondError(w, http.StatusInternalServerError, err)
	}
}
