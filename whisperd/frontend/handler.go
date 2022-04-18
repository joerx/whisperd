package frontend

import (
	"net/http"
)

type Opts struct{}

func New(opts Opts) (http.Handler, error) {
	fe := &frontend{}
	return fe, nil
}

type frontend struct{}

func (f *frontend) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	content := "<html><head><title>Hello World</title></head><body><h1>Hello World!</h1></body></html>"
	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(content))
}
