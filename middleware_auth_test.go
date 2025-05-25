// middleware_auth_test.go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"rss-scraper/internal/database"
)

func TestMiddlewareAuth(t *testing.T) {
	mockDB := &TestDB{}
	apiCfg := &testApiConfig{
		DB: mockDB,
	}

	t.Run("Valid API key", func(t *testing.T) {
		handlerCalled := false
		testHandler := func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlerCalled = true
			assert.Equal(t, "Test User", user.Name)
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "ApiKey valid-key")
		w := httptest.NewRecorder()

		middlewareFunc := apiCfg.middlewareAuth(testHandler)
		middlewareFunc(w, req)

		assert.True(t, handlerCalled)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid API key format", func(t *testing.T) {
		handlerCalled := false
		testHandler := func(w http.ResponseWriter, r *http.Request, user database.User) {
			handlerCalled = true
		}

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		middlewareFunc := apiCfg.middlewareAuth(testHandler)
		middlewareFunc(w, req)

		assert.False(t, handlerCalled)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
