// test_helpers.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"rss-scraper/internal/auth"
	"rss-scraper/internal/database"
	"time"
)

type DatabaseQuerier interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (database.User, error)
	GetFeeds(ctx context.Context) ([]database.Feed, error)
	CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error)
	CreateFeedFollow(ctx context.Context, arg database.CreateFeedFollowParams) (database.FeedFollow, error)
}

type testApiConfig struct {
	DB DatabaseQuerier
}

type TestDB struct{}

func (t *TestDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return database.User{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		Name:      arg.Name,
		ApiKey:    "test-api-key",
	}, nil
}

func (t *TestDB) GetUserByAPIKey(ctx context.Context, apiKey string) (database.User, error) {
	return database.User{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Test User",
		ApiKey:    apiKey,
	}, nil
}

func (t *TestDB) GetFeeds(ctx context.Context) ([]database.Feed, error) {
	return []database.Feed{
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "Test Feed",
			Url:       "https://example.com/feed",
			UserID:    uuid.New(),
		},
	}, nil
}

func (t *TestDB) CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.Feed, error) {
	return database.Feed{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		Name:      arg.Name,
		Url:       arg.Url,
		UserID:    arg.UserID,
	}, nil
}

func (t *TestDB) CreateFeedFollow(ctx context.Context, arg database.CreateFeedFollowParams) (database.FeedFollow, error) {
	return database.FeedFollow{
		ID:        arg.ID,
		CreatedAt: arg.CreatedAt,
		UpdatedAt: arg.UpdatedAt,
		UserID:    arg.UserID,
		FeedID:    arg.FeedID,
	}, nil
}

func (apiCfg *testApiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating user %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseToUser(user))
}

func (apiCfg *testApiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching feeds %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseFeedsToFeeds(feeds))
}

func (apiCfg *testApiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, http.StatusOK, databaseToUser(user))
}

func (apiCfg *testApiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseFeedToFeed(feed))
}

func (apiCfg *testApiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating feed follow %v", err))
		return
	}
	respondWithJson(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *testApiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusForbidden, fmt.Sprintf("Auth err: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
