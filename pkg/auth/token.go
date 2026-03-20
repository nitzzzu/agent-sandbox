package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/agent-sandbox/agent-sandbox/pkg/config"
)

func IsAuthEnabled() bool {
	return config.Cfg != nil && len(config.Cfg.APITokens) > 0
}

func ExtractToken(r *http.Request) string {
	token := strings.TrimSpace(r.Header.Get("X-Api-Key"))
	if token != "" {
		return token
	}
	return strings.TrimSpace(r.URL.Query().Get("api_key"))
}

func IsTokenAllowed(token string) bool {
	if !IsAuthEnabled() {
		return true
	}
	if token == "" {
		return false
	}
	for _, allowed := range config.Cfg.APITokens {
		if token == allowed {
			return true
		}
	}
	return false
}

func ValidateRequestToken(r *http.Request) (string, bool) {
	token := ExtractToken(r)
	if !IsTokenAllowed(token) {
		return "", false
	}

	return token, true
}

func GetUserTokenFromContext(ctx context.Context) string {
	value := ctx.Value("user")
	user, ok := value.(string)
	if !ok {
		return ""
	}
	return user
}
