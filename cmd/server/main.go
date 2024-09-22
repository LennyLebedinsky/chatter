package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lennylebedinsky/chatter/internal/domain"
	"github.com/lennylebedinsky/chatter/internal/gateway"
	"github.com/lennylebedinsky/chatter/internal/message"
)

type config struct {
	Host string
	Port string
}

func main() {
	config := &config{
		Host: "localhost",
		Port: "8080",
	}
	logger := log.Default()

	gw := gateway.New(
		domain.NewInMemoryRepository(),
		message.NewInMemoryStore(),
		logger)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: gw.Router(),
	}

	go func() {
		logger.Printf("HTTP server is listening on %s\n", httpServer.Addr+" ...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("error listening and serving: %s\n", err)
		}
	}()

	gw.StartBroadcaster(context.Background())

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting server down...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Printf("HTTP server forced to shutdown %v\n", err)
	}

	logger.Println("HTTP server had been shut down.")
}
