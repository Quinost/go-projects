package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Printf("Request for path: %s \n", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}