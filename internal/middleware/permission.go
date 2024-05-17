package middleware

import (
	"fmt"
	"github.com/wardrobe-gateway/internal/jwt"
	"net/http"
	"strings"
	"time"
)

const (
	googleValidationUrl = "https://www.googleapis.com/oauth2/v1/tokeninfo?access_token="
)

// просто проверяем что все ок с токеном

func WithJWTAuth(jwtService jwt.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(bearer, "Bearer ")
			if !ok {
				http.Error(w, "bearer token is absent", http.StatusUnauthorized)
				return
			}

			token, err := jwtService.Parse(tokenString)
			if err != nil {
				http.Error(w, fmt.Sprintf("cannot parsing token: %v", err), http.StatusUnauthorized)
				return
			}

			expirationTime, err := token.Claims.GetExpirationTime()
			if err != nil {
				http.Error(w, fmt.Sprintf("cannot parsing token: %v", err), http.StatusUnauthorized)
				return
			}

			if expirationTime.Time.Before(time.Now()) {
				http.Error(w, "token is expired", http.StatusUnauthorized)
				return
			}

			_, err = token.Claims.GetSubject()
			if err != nil {
				http.Error(w, fmt.Sprintf("cannot parsing token: %v", err), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func WithGoogleAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(bearer, "Bearer ")
			if !ok {
				http.Error(w, "Bearer token is absent", http.StatusUnauthorized)
				return
			}
			linkWithToken := fmt.Sprintf("%s%s", googleValidationUrl, tokenString)
			req, err := http.NewRequest("GET", linkWithToken, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
				return
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error making request to validation service: %v", err), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
