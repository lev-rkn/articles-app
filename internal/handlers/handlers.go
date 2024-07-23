package handlers

import (
	"fmt"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()
    mux.Handle("GET /", h("HELLO"))
	return mux
}

func h(name string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "%s: Вы вызвали %s методом %s\n", name, r.URL.String(), r.Method)
    }
}
