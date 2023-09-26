package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/smekuria1/podagg/internal/db"
	"github.com/smekuria1/podagg/models"
)

func (ApiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	feed, err := ApiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create feed: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DBFeedtoFeed(feed))
}

func (ApiCfg *ApiConfig) HandlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := ApiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't Get feeds: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DBFeedstoFeeds(feeds))
}
