package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	authorization string     = "Authorization"
	bearer        string     = "Bearer "
	claimsKey     contextKey = "claims"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/auth/")  {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get(authorization)
			if !strings.HasPrefix(authHeader, bearer) {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, bearer)
			token, err := jwt.Parse(tokenStr, keyFunc(secret))

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// claims := token.Claims.(jwt.MapClaims)
			// ctx := context.WithValue(r.Context(), claimsKey, claims)
			// next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
		})
	}
}

func keyFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}
