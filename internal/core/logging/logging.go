package logging

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger zerolog.Logger

func init() {
	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	writers = append(writers, &lumberjack.Logger{
		Filename:   "logs/livefetcher.log",
		MaxSize:    100, // megabytes
		MaxBackups: 100,
		MaxAge:     28, // days
	})
	mw := io.MultiWriter(writers...)
	logger = zerolog.New(mw).With().Stack().Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}
