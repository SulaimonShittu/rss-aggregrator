package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for handlerReadiness
func TestHandlerReadiness(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/healthz", nil)
	w := httptest.NewRecorder()

	handlerReadiness(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test for handlerErr
func TestHandlerErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/err", nil)
	w := httptest.NewRecorder()

	handlerErr(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp struct {
		Error string `json:"error"`
	}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Something went wrong", resp.Error)
}

// Test for handlerCreateUser
func TestHandlerCreateUser(t *testing.T) {
	mockDB := &TestDB{}
	apiCfg := &testApiConfig{
		DB: mockDB,
	}

	reqBody := map[string]string{
		"name": "Test User",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiCfg.handlerCreateUser(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Test User", resp["name"])
	assert.Equal(t, "test-api-key", resp["api_key"])
}

// Test for handlerGetFeeds
func TestHandlerGetFeeds(t *testing.T) {
	mockDB := &TestDB{}
	apiCfg := &testApiConfig{
		DB: mockDB,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/feeds", nil)
	w := httptest.NewRecorder()

	apiCfg.handlerGetFeeds(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp []map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "Test Feed", resp[0]["name"])
}
