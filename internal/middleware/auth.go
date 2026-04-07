package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	AccountID int
	Email     string
	Role      string
}

type authContextKey string

const authUserKey authContextKey = "auth-user"

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"missing or invalid authorization header"}`))
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid token"}`))
				return
			}

			mapClaims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid token claims"}`))
				return
			}

			accountIDFloat, ok := mapClaims["sub"].(float64)
			if !ok {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"invalid token subject"}`))
				return
			}

			email, _ := mapClaims["email"].(string)
			role, _ := mapClaims["role"].(string)

			claims := AuthClaims{
				AccountID: int(accountIDFloat),
				Email:     email,
				Role:      role,
			}

			ctx := context.WithValue(r.Context(), authUserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthClaims(ctx context.Context) (AuthClaims, bool) {
	claims, ok := ctx.Value(authUserKey).(AuthClaims)
	return claims, ok
}
