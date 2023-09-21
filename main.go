package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/smekuria1/podagg/handlers"
)

func main() {

	l := log.New(os.Stdout, "podagg-api", log.LstdFlags)

	sm := mux.NewRouter()
	v1Router := sm.PathPrefix("/v1").Methods(http.MethodGet).Subrouter()
	v1Router.HandleFunc("/healthz", handlers.HandlerReadiness)
	v1Router.HandleFunc("/err", handlers.HandlerErr)

	l.Printf("Server started on port: 8080")
	s := &http.Server{
		Addr:         "localhost:8080",
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
