package main

import (
	"context"
	"healthchecker/internal/checker"
	"healthchecker/internal/config"
	"healthchecker/internal/notifier"
	"healthchecker/internal/notifier/loggers"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	
	cfg, _ := config.Load("config.yaml")
	notifier := prepareLoggers()
	worker := checker.NewWorker(cfg, notifier)

	worker.Run(ctx)
}

func prepareLoggers() *notifier.Notifier {
	slogLogger := loggers.NewSlogLogger()
	return notifier.NewNotifier(slogLogger)
}
