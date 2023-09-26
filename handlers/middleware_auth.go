package handlers

import (
	"fmt"
	"net/http"

	"github.com/smekuria1/podagg/internal/auth"
	"github.com/smekuria1/podagg/internal/db"
)

type AuthHandler func(http.ResponseWriter, *http.Request, db.User)

func (cfg *ApiConfig) MiddleWareAuth(handler AuthHandler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		apikey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth erro: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apikey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)

	}
}
