package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/smekuria1/podagg/internal/db"
	"github.com/smekuria1/podagg/models"
)

func (ApiCfg *ApiConfig) HandlerFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	feedfollow, err := ApiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DBFeedFollowtoFeedFollow(feedfollow))
}

func (ApiCfg *ApiConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {

	feedfollow, err := ApiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DBFeedFollowstoFeedFollows(feedfollow))
}

func (ApiCfg *ApiConfig) HandlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {

	vars := mux.Vars(r)
	feedFollowsIDStr := vars["feedFollowID"]
	feedFollowId, err := uuid.Parse(feedFollowsIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't parse feed follow id: %v", err))
	}

	err = ApiCfg.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
