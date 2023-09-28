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

func (ApiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	user, err := ApiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, models.DBUserToUser(user))
}

func (ApiConfig *ApiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, 200, models.DBUserToUser(user))

}

func (ApiConfig *ApiConfig) HandleGetPostsForUser(w http.ResponseWriter, r *http.Request, user db.User) {
	posts, err := ApiConfig.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
	}

	respondWithJSON(w, 200, models.DBPostsToPosts(posts))
}
