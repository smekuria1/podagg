package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	l := log.New(os.Stdout, "podagg-api-json", log.LstdFlags)
	if code > 499 {
		l.Println("Responding with a 5XX error", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	l := log.New(os.Stdout, "podagg-api-json", log.LstdFlags)
	dat, err := json.Marshal(payload)
	if err != nil {
		l.Printf("Failed to marshal JSON response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)

}
