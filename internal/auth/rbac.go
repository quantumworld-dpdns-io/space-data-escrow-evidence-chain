package auth

import "net/http"

func IsAllowed(role, path, method string) bool {
	if role == "admin" {
		return true
	}
	if role == "operator" {
		if method == http.MethodPost && (path == "/v1/evidence" || path == "/v1/custody" || path == "/v1/attest" || path == "/v1/enrich") {
			return true
		}
		if method == http.MethodGet || method == http.MethodPost {
			return true
		}
	}
	if role == "viewer" {
		if method == http.MethodGet {
			return true
		}
		if method == http.MethodPost && (path == "/v1/verify" || path == "/v1/verify/bulk") {
			return true
		}
	}
	return false
}
