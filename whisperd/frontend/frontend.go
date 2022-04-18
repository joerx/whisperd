package frontend

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//go:embed static
var staticFS embed.FS

type Opts struct{}

func Handler() (http.Handler, error) {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", getFile)
	return mux, nil
}

func openFile(path string) ([]byte, error) {
	path = strings.TrimSuffix(path, "/")

	f, err := staticFS.Open(path)
	if err != nil {
		return []byte{}, err
	}

	fi, _ := f.Stat()
	if fi.IsDir() {
		return openFile(fmt.Sprintf("%s/index.html", path))
	}

	return ioutil.ReadAll(f)
}

func respondError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Add("Content-type", "text/plain")
	w.WriteHeader(statusCode)
	w.Write([]byte("not found"))
}

func sendFile(w http.ResponseWriter, data []byte) {
	ct := http.DetectContentType(data)
	w.Header().Add("Content-type", ct)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getFile(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	log.Printf("GET %s", path)

	data, err := openFile(fmt.Sprintf("static%s", path))

	switch {
	case err == nil:
		sendFile(w, data)
	case errors.Is(err, fs.ErrNotExist):
		respondError(w, http.StatusNotFound, "not found")
	default:
		log.Println(err)
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
