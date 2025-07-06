package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test: POST /login with missing username and password returns BadRequest
func TestAuthenticateUser_MissingFields(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", strings.NewReader("username=&password="))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	authenticateUser(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Test: GET /login returns MethodNotAllowed
func TestAuthenticateUser_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	authenticateUser(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

// Test: GET /view_cert returns a response (status code may depend on your logic)
func TestViewCertificateForm(t *testing.T) {
	req := httptest.NewRequest("GET", "/view_cert", nil)
	w := httptest.NewRecorder()

	ViewCertificateForm(w, req)

	resp := w.Result()
	assert.NotNil(t, resp)
	// Optionally check for a specific status code or content
}

// Example: Test checkSession middleware structure
func TestCheckSessionMiddleware(t *testing.T) {
	handler := checkSession()(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	assert.NotNil(t, resp)
	// Optionally check for a specific status code
}
