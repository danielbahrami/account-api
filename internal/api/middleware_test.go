package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Missing Authorization header returns 401
func TestAuthMiddleware_NoToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rr := httptest.NewRecorder()

	authMiddleware(func(w http.ResponseWriter, r *http.Request) {})(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", rr.Code)
	}
}

// Invalid JWT returns 401
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "Secret")

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req.Header.Set("Authorization", "Bearer InvalidToken")

	rr := httptest.NewRecorder()

	authMiddleware(func(w http.ResponseWriter, r *http.Request) {})(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized, got %d", rr.Code)
	}
}
