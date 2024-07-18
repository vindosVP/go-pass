// Package sl initializes and configures slog logger
package sl

import (
	"log/slog"
	"os"

	"github.com/vindosVP/go-pass/pkg/logger/handlers/slogdiscard"
	"github.com/vindosVP/go-pass/pkg/logger/handlers/slogpretty"
)

const (
	// envLocal is a local environment
	envLocal = "local"

	// envLocal is a development environment
	envDev = "dev"

	// envLocal is a production environment
	envProd = "prod"

	// envTest is a test environment
	envTest = "test"
)

// Log consists configured logger instance
var Log *slog.Logger

// SetupLogger configures the logger depending on environment
func SetupLogger(env string) {
	var log *slog.Logger
	switch env {
	case envLocal:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}

		handler := opts.NewPrettyHandler(os.Stdout)

		log = slog.New(handler)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envTest:
		log = slog.New(
			slogdiscard.NewDiscardHandler(),
		)
	}
	Log = log
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
