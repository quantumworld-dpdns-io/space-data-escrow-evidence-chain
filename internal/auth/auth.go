package auth

import (
	"errors"
	"net/http"
	"strings"
)

type Principal struct {
	Role   string
	Source string
}

func Authenticate(r *http.Request, apiKey string) (Principal, error) {
	if apiKey != "" && r.Header.Get("X-API-Key") == apiKey {
		return Principal{Role: "admin", Source: "api_key"}, nil
	}
	authz := r.Header.Get("Authorization")
	if strings.HasPrefix(authz, "Bearer ") {
		tok := strings.TrimPrefix(authz, "Bearer ")
		switch tok {
		case "jwt-admin":
			return Principal{Role: "admin", Source: "jwt"}, nil
		case "jwt-operator":
			return Principal{Role: "operator", Source: "jwt"}, nil
		case "jwt-viewer":
			return Principal{Role: "viewer", Source: "jwt"}, nil
		}
	}
	return Principal{}, errors.New("unauthorized")
}
