package log

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	log "github.com/thedataflows/go-lib-log"
)

func Api(r *http.Request) {
	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	logger.Debug().Msg("(API) \"" + r.Pattern + "\"")
}

func Database(action string, err error) {
	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	logger.Debug().AnErr(" (DATABASE) | \""+action+"\"", err)
}

func Info(info string) {
	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.InfoLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	logger.Info().Msg(info)
}

func Fatal(error error) {
	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.FatalLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	logger.Fatal().Err(error)
}
