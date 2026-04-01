package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/codingzen/living-circles-mrl/src/server/transport"
)

func main() {
	server := transport.NewServer()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := server.Run(ctx); err != nil {
			log.Printf("server loop stopped: %v", err)
		}
	}()

	log.Println("listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", server.Handler()); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
