package notifier

type StatusChanged struct {
	URL    string
	CurrentStatus string
	PreviousStatus string
}

type Logger interface {
	NotifyChanged(message StatusChanged)
	Log(message string)
}

type Notifier struct {
	loggers []Logger
}

func NewNotifier(loggers ...Logger) *Notifier {
	return &Notifier{loggers: loggers}
}

func (l *Notifier) Notify(message StatusChanged) {
	for _, logger := range l.loggers {
		logger.NotifyChanged(message)
	}
}

func (l *Notifier) Log(message string) {
	for _, logger := range l.loggers {
		logger.Log(message)
	}
}

