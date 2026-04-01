package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codingzen/living-circles-mrl/src/server/transport"
)

func main() {
	app := transport.NewServer()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := app.Run(ctx); err != nil {
			log.Printf("server loop stopped: %v", err)
		}
	}()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: app.Handler(),
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("http shutdown failed: %v", err)
		}
	}()

	log.Println("listening on http://localhost:8080")
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
