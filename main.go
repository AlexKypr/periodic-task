package main

import (
	"context"
	"inaccess/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	envPort     = "APP_PORT"
	defaultPort = "8080"
)

// NewServer helper function to inject logger and if everything else that it may be added in the future at handlers
func NewServer(l *log.Logger) *http.Server {
	pl := handlers.NewPtlistQuery(l)

	sm := http.NewServeMux()
	sm.Handle("/ptlist", pl)

	address := ":" + os.Getenv(envPort)
	return &http.Server{
		Addr:         address,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func main() {
	if _, ok := os.LookupEnv(envPort); !ok {
		os.Setenv(envPort, defaultPort)
	}

	// create "custom" logger to pass it on our api
	l := log.New(os.Stdout, "ts-api ", log.LstdFlags)

	server := NewServer(l)

	go func() {
		l.Println("Starting server on port ", os.Getenv(envPort))

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// notify in case user press ctrl + c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block main goroutine until we receive a termination signal
	sig := <-c
	l.Println("got termination signal: ", sig)

	// create context that will cancel in 30 seconds, keeping connection alive for 30 seconds to fulfill their tasks
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shutdown server gracefully
	err := server.Shutdown(ctx)
	if err != nil {
		l.Fatal(err)
	}
}
