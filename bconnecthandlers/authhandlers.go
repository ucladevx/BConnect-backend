package bconnecthandlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/ucladevx/BConnect-backend/storage/postgres"
)

type key int

const (
	keyPrincipalID key = iota
)

// VerifyToken Verifies auth token
func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("x-access-token") //Grab the token from the header

		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			return
		}
		tk := &postgres.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), keyPrincipalID, tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
