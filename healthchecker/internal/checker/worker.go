package checker

import (
	"context"
	"healthchecker/internal/config"
	"healthchecker/internal/notifier"
	"net/http"
	"sync"
)

type Worker struct {
	monitor  *[]Monitor
	config   *config.Config
	notifier *notifier.Notifier
	wg       *sync.WaitGroup
}

func NewWorker(cfg *config.Config, notifier *notifier.Notifier) *Worker {
	return &Worker{
		config:   cfg,
		notifier: notifier,
		wg:       &sync.WaitGroup{},
	}
}

func (w *Worker) Run(ctx context.Context) {

	httpClient := &http.Client{
		Timeout: w.config.Settings.RequestTimeout,
	}

	for _, site := range w.config.Pages {
		w.wg.Add(1)
		monitor := NewMonitor(&site, httpClient, w.wg, w.notifier)
		go monitor.StartMonitoring(ctx)
	}

	<-ctx.Done()
	w.notifier.Log("Shutting down worker...")
	w.wg.Wait()
	w.notifier.Log("Worker shut down.")
}
