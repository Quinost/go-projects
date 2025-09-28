package middleware

import (
	"context"
	"net/http"
	"strings"
	"warehouse/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string, anonymousPrefixes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAnonymousPrefix(r.URL.Path, anonymousPrefixes) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get(constants.Authorization)
			if !strings.HasPrefix(authHeader, constants.Bearer) {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, constants.Bearer)
			token, err := jwt.Parse(tokenStr, keyFunc(secret))

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			ctx := context.WithValue(r.Context(), constants.ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func keyFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func isAnonymousPrefix(path string, allowedPrefixes []string) bool {
    for _, prefix := range allowedPrefixes {
        if strings.HasPrefix(path, prefix) {
            return true
        }
    }
    return false
}
