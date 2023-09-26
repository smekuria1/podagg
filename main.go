package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/smekuria1/podagg/handlers"
	"github.com/smekuria1/podagg/internal/db"
	"github.com/spf13/viper"
)

func main() {

	l := log.New(os.Stdout, "podagg-api", log.LstdFlags)
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()
	if err != nil {
		l.Fatal("Error Reading config file")
	}
	portString := viper.GetString("PORT")
	dbURL := viper.GetString("DB_URL")

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		l.Fatal("Couldn't establish connection to DB")
	}

	apiCfg := handlers.ApiConfig{
		DB: db.New(conn),
	}

	sm := mux.NewRouter()
	v1RouterGET := sm.PathPrefix("/v1").Methods(http.MethodGet).Subrouter()
	v1RouterGET.HandleFunc("/healthz", handlers.HandlerReadiness)
	v1RouterGET.HandleFunc("/err", handlers.HandlerErr)
	v1RouterGET.HandleFunc("/users", apiCfg.MiddleWareAuth(apiCfg.HandleGetUser))
	v1RouterGET.HandleFunc("/feeds", apiCfg.HandlerGetFeeds)

	v1RouterPOST := sm.PathPrefix("/v1").Methods(http.MethodPost).Subrouter()
	v1RouterPOST.HandleFunc("/users", apiCfg.HandlerCreateUser)
	v1RouterPOST.HandleFunc("/feeds", apiCfg.MiddleWareAuth(apiCfg.HandlerCreateFeed))

	l.Printf("Server started on port: %s", portString)
	s := &http.Server{
		Addr:         "localhost:" + portString,
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	sig := <-sigChan

	l.Println("Recived Terminate, Shutting Down Gracefully", sig)
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(tc)

}
