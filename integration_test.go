// integration_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Setup a router with our test version for integration tests
func setupTestRouter() *chi.Mux {
	mockDB := &TestDB{}
	apiCfg := &testApiConfig{
		DB: mockDB,
	}

	router := chi.NewRouter()
	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

	router.Mount("/v1", v1Router)
	return router
}

// Test the health endpoint
func TestIntegrationHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/v1/healthz", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test creating a user and then getting the user
func TestIntegrationUserFlow(t *testing.T) {
	router := setupTestRouter()

	userBody := map[string]string{
		"name": "Integration User",
	}

	userJSON, _ := json.Marshal(userBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NoError(t, err)
	assert.Equal(t, "Integration User", createResp["name"])
	assert.Equal(t, "test-api-key", createResp["api_key"])

	req = httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	req.Header.Set("Authorization", "ApiKey test-api-key")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getUserResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &getUserResp)
	assert.NoError(t, err)
	assert.Equal(t, "Integration User", getUserResp["name"])
}

// Test feeds endpoint
func TestIntegrationGetFeeds(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/v1/feeds", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp []map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Test Feed", resp[0]["name"])
}
