package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"rssagg/internal/database"
)

func (apiconf *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	user, err := apiconf.DB.CreateUser(
		r.Context(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      params.Name,
		},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
	}
	respondWithJSON(w, 201, databaseUsertoUser(user))
}

func (apiconf *apiConfig) handlerGetUser(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	respondWithJSON(w, 200, databaseUsertoUser(user))
}

func (apiconf *apiConfig) handlerGetPostsForUser(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	posts, err := apiconf.DB.GetPostsForUser(
		r.Context(),
		database.GetPostsForUserParams{UserID: user.ID, Limit: 10},
	)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
	}
	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
