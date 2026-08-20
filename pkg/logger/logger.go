package logger

import (
	"log/slog"
	"os"
)

// New returns a JSON slog.Logger with the time and message keys renamed to
// "timestamp" and "message".
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				a.Key = "timestamp"
			case slog.MessageKey:
				a.Key = "message"
			}
			return a
		},
	})
	return slog.New(handler)
}
