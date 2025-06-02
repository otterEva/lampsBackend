package logs

import (
	"log/slog"
	"os"
	"time"

	"github.com/otterEva/lamps/orders_service/settings"
)

func initLogger() *slog.Logger {

	levelStr := settings.Config.LOG_LEVEL

	var level slog.Level

	switch levelStr {
	case "Debug":
		level = slog.LevelDebug
	case "Info":
		level = slog.LevelInfo
	case "Warn":
		level = slog.LevelWarn
	case "Error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Any().(time.Time)
				formatted := t.Format("15:04:05")
				return slog.Attr{Key: a.Key, Value: slog.StringValue(formatted)}
			}
			return a

		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
	return logger
}

var Logger = initLogger()
