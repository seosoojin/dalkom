package middlewares

import (
	"context"
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
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			if !strings.HasPrefix(bearer, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			token := strings.TrimPrefix("Bearer ", bearer)

			_, err := a.JWTService.VerifyToken(token)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), auth.CTXJWTKEY, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
