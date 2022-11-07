package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"periodictask/handlers"
	"periodictask/utils"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// NewServer helper function to inject logger and if everything else that it may be added in the future at handlers
func NewServer(l *zap.Logger, args *utils.Args) *http.Server {
	th := handlers.NewTaskHandler(l)

	sm := http.NewServeMux()
	sm.Handle("/ptlist", th)

	address := fmt.Sprintf(":%s", args.Port)
	return &http.Server{
		Addr:         address,
		Handler:      sm,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

func main() {
	// create "custom" logger to pass it on our api
	l, err := zap.NewProduction(zap.AddCaller())
	if err != nil {
		fmt.Println("Error creating logger: ", err)
		os.Exit(1)
	}

	args, err := utils.ParseArgs()
	if err != nil {
		l.Sugar().Fatalf("Error while parsing args: %s", err)
	}

	errC := run(l, args)
	if err = <-errC; err != nil {
		l.Sugar().Fatalf("Error while server was running: %s", err)
	}
}

func run(l *zap.Logger, args *utils.Args) <-chan error {
	srv := NewServer(l, args)
	errC := make(chan error, 1)

	go func() {
		l.Sugar().Infof("Server is listening on %s:%s", args.Address, args.Port)

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

		// block goroutine until we receive a termination signal
		sig := <-c
		l.Sugar().Infof("Shutdown signal received: ", sig)

		// create context that will cancel in 5 seconds, keeping connection alive for 5 seconds to fulfill their tasks
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			l.Sync()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		// shutdown server gracefully
		if err := srv.Shutdown(ctx); err != nil {
			errC <- err
		}
		l.Info("Shutdown completed")
	}()

	return errC
}
