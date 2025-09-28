package middleware

import (
	"log"
	"net/http"
	"warehouse/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user_name := ""
		if claims, ok := r.Context().Value(constants.ClaimsKey).(jwt.MapClaims); ok {
			if user, ok := claims[constants.ClaimUserName].(string); ok {
				user_name = user
			}
		}

		log.Printf("Request path: %s user_name: %s \n", r.URL.Path, user_name)
		next.ServeHTTP(w, r)
	})
}
