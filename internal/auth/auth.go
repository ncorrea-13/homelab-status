package auth

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

type Auth struct {
	secret string
}

func NewAuth(secret string) *Auth {
	return &Auth{secret: secret}
}

func LoadSecret() (string, error) {
	filePath := os.Getenv("WEBHOOK_SECRET_FILE")
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(content)), nil
	}
	return os.Getenv("WEBHOOK_SECRET"), nil
}

func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provided := r.Header.Get("X-Webhook-Secret")

		if provided == "" {
			provided = r.URL.Query().Get("secret")
		}

		if a.secret == "" || provided == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if subtle.ConstantTimeCompare([]byte(provided), []byte(a.secret)) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
