package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"thumbnail/internal/api"
	"thumbnail/internal/services/thumbnail"
	"thumbnail/internal/services/worker"
)

func main() {
	fmt.Println("Starting server...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := worker.NewPool[thumbnail.ThumbnailJob](2, 100)
	pool.Start(ctx)
	thumbnailService := thumbnail.NewThumbnailService(pool)
	api, serverErrChan := runApi(thumbnailService)

	handleShutdown(ctx, serverErrChan, api, pool)
}

func runApi(thumbnailService *thumbnail.ThumbnailService) (*api.Server, chan error) {
	api := api.NewServer(thumbnailService)

	serverErrChan := make(chan error, 1)
	go func() {
		if err := api.Start(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	return api, serverErrChan
}

func handleShutdown(ctx context.Context, serverErrChan <-chan error, api *api.Server, pool *worker.Pool[thumbnail.ThumbnailJob]) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-shutdownCh:
		fmt.Print("Shutdown...\n", sig)
	case err := <-serverErrChan:
		fmt.Printf("Server error: %v\n", err)
	}

	api.Stop(ctx)
	pool.Stop()

	fmt.Println("Stopped")
}
