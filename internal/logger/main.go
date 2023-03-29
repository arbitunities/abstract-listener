package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func NewLogger(debug *bool) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug || (strings.ToUpper(os.Getenv("LOG_LEVEL")) == "DEBUG") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return zerolog.New(os.Stdout)
}
