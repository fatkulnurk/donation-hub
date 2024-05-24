package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// BearerTokenMiddleware middleware for checking Bearer token
func authTokenMiddleware(next http.HandlerFunc, api *API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(NewInvalidAccessToken())
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		fmt.Println("Token:", token)
		ctx := context.WithValue(r.Context(), "token", token)
		fmt.Println(token)

		payload, err := api.AuthToken.ValidateToken(token)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(NewInvalidAccessToken())
			return
		}

		ctx = context.WithValue(ctx, "auth_id", payload.UserID)
		ctx = context.WithValue(ctx, "auth_username", payload.Username)
		ctx = context.WithValue(ctx, "auth_email", payload.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func toJsonMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
