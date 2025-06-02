package api

import (
	"fmt"
	"net/http"
)

func (a *Api) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.write(fmt.Sprintf("%s %s %s %s\n", r.Proto, r.Method, r.RequestURI, r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}

func (a *Api) contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
