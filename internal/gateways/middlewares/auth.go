package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/seosoojin/dalkom/internal/domain/auth"
)

type Authenticator interface {
	Authenticate() func(next http.Handler) http.Handler
}

type authenticator struct {
	JWTService auth.JWTService
}

var _ Authenticator = &authenticator{}

func NewAuthenticator(jwtService auth.JWTService) *authenticator {
	return &authenticator{
		JWTService: jwtService,
	}
}

func (a *authenticator) Authenticate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get("Authorization")
			if bearer == "" {
				log.Println("invalid token")
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			if !strings.HasPrefix(bearer, "Bearer ") {
				log.Println("invalid token")
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			token := strings.TrimPrefix(bearer, "Bearer ")

			jwt, err := a.JWTService.VerifyToken(token)
			if err != nil {
				log.Println(err)
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), auth.CTXJWTKEY, jwt.Claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
