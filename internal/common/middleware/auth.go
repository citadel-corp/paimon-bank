package middleware

import (
	"context"
	"net/http"

	"github.com/citadel-corp/paimon-bank/internal/common/jwt"
)

type ContextAuthKey struct{}

func Authorized(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(tokenString) <= len("Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[len("Bearer "):]
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		subject, err := jwt.VerifyAndGetSubject(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextAuthKey{}, subject)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// Authenticate request only if authorization header is set
func Authenticate(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			next(w, r)
			return
		}

		if len(tokenString) <= len("Bearer ") {
			next(w, r)
			return
		}

		tokenString = tokenString[len("Bearer "):]
		if tokenString == "" {
			next(w, r)
			return
		}

		subject, err := jwt.VerifyAndGetSubject(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextAuthKey{}, subject)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
