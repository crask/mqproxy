package server

import (
	"io"
	"net/http"
)

// Http Handler
type HttpHandler struct {
	Mux map[string]func(http.ResponseWriter, *http.Request)
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	serveHTTP(h.Mux, w, r)
}

func serveHTTP(mux map[string]func(http.ResponseWriter, *http.Request), w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.Path]; ok {
		h(w, r)
		return
	}

	//TODO: handle 404, just for debugging now
	io.WriteString(w, "My server: "+r.URL.Path)
}
