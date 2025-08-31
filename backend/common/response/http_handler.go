package response

import (
	"fmt"
	"net/http"
)

type NotFoundHandler struct {
}

func (h NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	NotFound(w, fmt.Sprintf("Method %s %s not found", r.Method, r.URL.Path))
}

type NotAllowedHandler struct {
}

func (h NotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	NotFound(w, fmt.Sprintf("Method %s %s not allowed", r.Method, r.URL.Path))
}
