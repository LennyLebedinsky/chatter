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

	"github.com/lennylebedinsky/chatter/internal/chat"
	"github.com/lennylebedinsky/chatter/internal/domain"
	"github.com/lennylebedinsky/chatter/internal/gateway"
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
	repo := domain.NewInMemoryRepository()
	gw := gateway.New(
		chat.NewBroadcaster(repo, logger),
		repo,
		logger)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: gw.Router(),
	}

	//stopBroadcast := make(chan struct{})

	go func() {
		logger.Printf("HTTP server is listening on %s\n", httpServer.Addr+" ...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("error listening and serving: %s\n", err)
		}
	}()

	go func() {
		logger.Println("Starting message broadcaster...")
		gw.Broadcaster().Start()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting server down...")

	//close(stopBroadcast)

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Printf("HTTP server forced to shutdown %v\n", err)
	}

	logger.Println("HTTP server had been shut down.")
}
