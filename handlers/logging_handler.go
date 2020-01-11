package handlers

import (
	"net/http"
	"log"
)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r:=recover(); r != nil {
				http.Error(w, "Error code 1", http.StatusInternalServerError)
			}
		}()
		log.Printf("%s\t %s\t %s\t %s\t", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

		next.ServeHTTP(w, r)
	})
}