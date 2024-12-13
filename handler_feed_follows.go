package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"rssagg/internal/database"
)

func (apiconf *apiConfig) handlerCreateFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	feedFollow, err := apiconf.DB.CreateFeedFollow(
		r.Context(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    params.FeedID,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %v", err))
	}

	respondWithJSON(w, 201, databaseFeedFollowtoFeedFollow(feedFollow))
}

func (apiconf *apiConfig) handlerGetFeedFollows(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedFollows, err := apiconf.DB.GetFeedFollows(
		r.Context(),
		user.ID,
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follows: %v", err))
	}

	respondWithJSON(w, 201, databaseFeedFollowstoFeedFollows(feedFollows))
}

func (apiconf *apiConfig) handlerDeleteFeedFollow(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow ID: %v", err))
	}
	err = apiconf.DB.DeleteFeedFollows(
		r.Context(),
		database.DeleteFeedFollowsParams{ID: feedFollowID, UserID: user.ID},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follows: %v", err))
	}

	respondWithJSON(w, 200, struct{}{})
}
