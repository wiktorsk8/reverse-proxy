package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wiktorsk8/reverse-proxy/internal/config"
)

type JWTAuth struct {
	Config config.AuthConfig
}

func NewJWTAuthMiddleware(config config.AuthConfig) *JWTAuth {
	return &JWTAuth{
		Config: config,
	}
}

func (j *JWTAuth) GetMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := j.Config.JWTSecret

			tokenString, err := j.retrieveTokenString(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secret), nil
			})

			if err != nil {
				fmt.Println(err)
				http.Error(w, "JWT parsing failed.", http.StatusBadGateway)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)

			if !ok || !token.Valid {
				http.Error(w, "JWT token is invalid.", http.StatusUnauthorized)
				return
			}

			issuer, err := claims.GetIssuer()
			if err != nil {
				http.Error(w, "unexpected error", http.StatusBadGateway)
				return
			}

			r.Header.Del("X-User-Id")
			r.Header.Set("X-User-Email", issuer)

			next.ServeHTTP(w, r)
		})
	}
}

func (j *JWTAuth) retrieveTokenString(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("authorization header is required")
	}
	return strings.TrimPrefix(authorizationHeader, "Bearer "), nil
}
