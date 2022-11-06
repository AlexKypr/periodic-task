package main

import (
	"context"
	"fmt"
	"inaccess/handlers"
	"inaccess/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// NewServer helper function to inject logger and if everything else that it may be added in the future at handlers
func NewServer(l *log.Logger, args *utils.Args) *http.Server {
	pl := handlers.NewPtlistQuery(l)

	sm := http.NewServeMux()
	sm.Handle("/ptlist", pl)

	address := fmt.Sprintf(":%s", args.Port)
	fmt.Println(address)
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
	// create "custom" logger to pass it on our api
	l := log.New(os.Stdout, "ts-api ", log.LstdFlags)

	args, err := utils.ParseArgs()
	if err != nil {
		l.Fatalf("Error while parsing args: %s\n", err)
	}

	errC := run(l, args)
	if err = <-errC; err != nil {
		l.Fatalf("Error while server was running: %s\n", err)
	}
}

func run(l *log.Logger, args *utils.Args) <-chan error {

	srv := NewServer(l, args)

	errC := make(chan error, 1)

	go func() {
		l.Printf("Server is listening on %s:%s\n", args.Address, args.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	go func() {
		// notify in case user press ctrl + c
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		// block main goroutine until we receive a termination signal
		sig := <-c
		l.Println("\nShutdown signal received: ", sig)

		// create context that will cancel in 5 seconds, keeping connection alive for 5 seconds to fulfill their tasks
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		// shutdown server gracefully
		if err := srv.Shutdown(ctx); err != nil {
			errC <- err
		}
		l.Println("Shutdown completed")
	}()

	return errC
}
