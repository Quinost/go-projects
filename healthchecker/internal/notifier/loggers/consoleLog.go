package loggers

import (
	"healthchecker/internal/notifier"
	"log/slog"
	"os"
)

type SlogLog struct{
	logger *slog.Logger
}

func NewSlogLogger() *SlogLog {
	return &SlogLog{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (c *SlogLog) NotifyChanged(m notifier.StatusChanged) {
	c.logger.Warn("Status change detected for " + m.URL + ": " + string(m.PreviousStatus) + " -> " + string(m.CurrentStatus))
}

func (c *SlogLog) Log(message string) {
	c.logger.Info(message)
}
