package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"rssagg/internal/database"
)

func (apiconf *apiConfig) handlerCreateFeed(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	feed, err := apiconf.DB.CreateFeed(
		r.Context(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
			Url:       params.URL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %v", err))
	}

	respondWithJSON(w, 201, databaseFeedtoFeed(feed))
}

func (apiconf *apiConfig) handlerGetFeeds(
	w http.ResponseWriter,
	r *http.Request,
) {
	feeds, err := apiconf.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get feeds: %v", err))
	}
	respondWithJSON(w, 200, databaseFeedstoFeeds(feeds))
}