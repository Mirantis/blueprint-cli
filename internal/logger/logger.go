package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewLogger creates a new logger with the given log level
func NewLogger(logLevel string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, PartsExclude: []string{zerolog.TimestampFieldName}})
	zerolog.SetGlobalLevel(parseLevel(logLevel))
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
