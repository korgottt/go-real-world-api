package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"realworldapi/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGETUsers(t *testing.T) {
	server := NewGlobalServer()

	t.Run("returns user struct", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := defaultUserInfo

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns 200 on /api/user", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestPOSTuser(t *testing.T) {
	server := NewGlobalServer()

	t.Run("Authentication test", func(t *testing.T) {
		requestBody := `{
			"user":{
			  "email": "jake@jake.jake",
			  "password": "jakejake"
			}
		  }`

		request, _ := http.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(requestBody))
		request.Header.Set("content-type", "application/json; charset=utf-8")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		var got model.SingleUserWrap
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", response.Body, err)
		}

		assert.Equal(t, "jake@jake.jake", got.User.Email)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Registration test", func(t *testing.T) {
		requestBody := `{
			"user":{
			  "username": "Jacob",
			  "email": "jake@jake.jake",
			  "password": "jakejake"
			}
		  }`

		request, _ := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(requestBody))
		request.Header.Set("content-type", "application/json; charset=utf-8")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		var got model.SingleUserWrap
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from %q into User struct, '%v'", response.Body, err)
		}

		assert.Equal(t, "jakejake", got.User.Password)
		assertStatus(t, response.Code, http.StatusOK)
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
