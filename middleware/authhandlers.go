package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type key int

const (
	keyPrincipalID key = iota
)

// VerifyToken Verifies auth token
func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := strings.TrimSpace(r.Header.Get("x-access-token"))

		if header == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		type Claims struct {
			jwt.MapClaims
			UUID           string
			FirstName      string
			Email          string
			StandardClaims *jwt.StandardClaims
		}
		claims := Claims{}

		header = strings.Replace(header, "Bearer ", "", -1)
		_, err := jwt.ParseWithClaims(header, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if time.Now().Unix()-claims.StandardClaims.ExpiresAt > 15 {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), keyPrincipalID, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
