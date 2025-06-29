package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mbilaljawwad/trendy-repos/internal/config"
	"github.com/mbilaljawwad/trendy-repos/internal/server"
)

const (
	PORT = ":8080"
)

func main() {
	fmt.Println("Starting the application...")
	config.InitConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	appServer := server.NewAppServer(ctx)

	go func() {
		appServer.Start()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("Received signal to terminate the server")
	case <-ctx.Done():
		fmt.Println("Context cancel, initiating graceful shutdown")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	if err := appServer.Server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	fmt.Println("Server shutdown complete")

}
