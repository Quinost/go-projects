package checker

import (
	"context"
	"healthchecker/internal/config"
	"healthchecker/internal/notifier"
	"net/http"
	"sync"
	"time"
)

type Status string

const (
	StatusUP   Status = "UP"
	StatusDOWN Status = "DOWN"
)

type Monitor struct {
	site       *config.Site
	httpClient *http.Client
	wg         *sync.WaitGroup
	notifier   *notifier.Notifier
	lastStatus Status
}

func NewMonitor(site *config.Site, httpClient *http.Client, wg *sync.WaitGroup, notifier *notifier.Notifier) *Monitor {
	return &Monitor{
		site:       site,
		httpClient: httpClient,
		wg:         wg,
		notifier:   notifier,
		lastStatus: StatusUP,
	}
}

func (m *Monitor) StartMonitoring(ctx context.Context) {
	m.notifier.Log("Starting monitoring for " + m.site.URL)
	ticker := time.NewTicker(m.site.Interval)
	defer ticker.Stop()
	defer m.wg.Done()

	m.checkSite()

	for {
		select {
		case <-ticker.C:
			m.checkSite()
		case <-ctx.Done():
			m.notifier.Log("Stopping monitoring for " + m.site.URL)
			return
		}
	}
}

func (m *Monitor) checkSite() {
	response, err := m.httpClient.Get(m.site.URL)

	currentStatus := StatusUP

	if err != nil {
		currentStatus = StatusDOWN
	} else {
		if response.StatusCode < 200 || response.StatusCode >= 300 {
			currentStatus = StatusDOWN
		}
		response.Body.Close()
	}

	if currentStatus != m.lastStatus {
		statusChanged := notifier.StatusChanged{
			URL:            m.site.URL,
			CurrentStatus:  string(currentStatus),
			PreviousStatus: string(m.lastStatus),
		}
		m.notifier.Notify(statusChanged)
		m.lastStatus = currentStatus
	}
}
