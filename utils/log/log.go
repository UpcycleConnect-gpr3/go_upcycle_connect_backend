package log

import (
	"encoding/json"
	"go-upcycle_connect-backend/utils/rules"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	log "github.com/thedataflows/go-lib-log"
)

var httpStatusTexts = map[int]string{
	http.StatusContinue:                      "Continue",
	http.StatusSwitchingProtocols:            "Switching Protocols",
	http.StatusProcessing:                    "Processing",
	http.StatusOK:                            "OK",
	http.StatusCreated:                       "Created",
	http.StatusAccepted:                      "Accepted",
	http.StatusNonAuthoritativeInfo:          "Non-Authoritative Information",
	http.StatusNoContent:                     "No Content",
	http.StatusResetContent:                  "Reset Content",
	http.StatusPartialContent:                "Partial Content",
	http.StatusMultiStatus:                   "Multi-Status",
	http.StatusAlreadyReported:               "Already Reported",
	http.StatusIMUsed:                        "IM Used",
	http.StatusMultipleChoices:               "Multiple Choices",
	http.StatusMovedPermanently:              "Moved Permanently",
	http.StatusFound:                         "Found",
	http.StatusSeeOther:                      "See Other",
	http.StatusNotModified:                   "Not Modified",
	http.StatusUseProxy:                      "Use Proxy",
	http.StatusTemporaryRedirect:             "Temporary Redirect",
	http.StatusPermanentRedirect:             "Permanent Redirect",
	http.StatusBadRequest:                    "Bad Request",
	http.StatusUnauthorized:                  "Unauthorized",
	http.StatusPaymentRequired:               "Payment Required",
	http.StatusForbidden:                     "Forbidden",
	http.StatusNotFound:                      "Not Found",
	http.StatusMethodNotAllowed:              "Method Not Allowed",
	http.StatusNotAcceptable:                 "Not Acceptable",
	http.StatusProxyAuthRequired:             "Proxy Authentication Required",
	http.StatusRequestTimeout:                "Request Timeout",
	http.StatusConflict:                      "Conflict",
	http.StatusGone:                          "Gone",
	http.StatusLengthRequired:                "Length Required",
	http.StatusPreconditionFailed:            "Precondition Failed",
	http.StatusRequestEntityTooLarge:         "Request Entity Too Large",
	http.StatusRequestURITooLong:             "Request URI Too Long",
	http.StatusUnsupportedMediaType:          "Unsupported Media Type",
	http.StatusRequestedRangeNotSatisfiable:  "Requested Range Not Satisfiable",
	http.StatusExpectationFailed:             "Expectation Failed",
	http.StatusTeapot:                        "I'm a teapot",
	http.StatusUnprocessableEntity:           "Unprocessable Entity",
	http.StatusLocked:                        "Locked",
	http.StatusFailedDependency:              "Failed Dependency",
	http.StatusTooEarly:                      "Too Early",
	http.StatusUpgradeRequired:               "Upgrade Required",
	http.StatusPreconditionRequired:          "Precondition Required",
	http.StatusTooManyRequests:               "Too Many Requests",
	http.StatusRequestHeaderFieldsTooLarge:   "Request Header Fields Too Large",
	http.StatusUnavailableForLegalReasons:    "Unavailable For Legal Reasons",
	http.StatusInternalServerError:           "Internal Server Error",
	http.StatusNotImplemented:                "Not Implemented",
	http.StatusBadGateway:                    "Bad Gateway",
	http.StatusServiceUnavailable:            "Service Unavailable",
	http.StatusGatewayTimeout:                "Gateway Timeout",
	http.StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	http.StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	http.StatusInsufficientStorage:           "Insufficient Storage",
	http.StatusLoopDetected:                  "Loop Detected",
	http.StatusNotExtended:                   "Not Extended",
	http.StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

const (
	ErrAuthFailed    = "Credentials do not match"
	ErrInvalidBody   = "Incorrect body format"
	ErrJson          = "Json parse error"
	ErrGenerateToken = "Generate token error"
)

func Api(r *http.Request) {
	logger := log.NewLoggerBuilder().WithLogLevel(zerolog.DebugLevel).WithBufferSize(10000).WithRateLimit(1000).WithGroupWindow(2 * time.Second).WithOutput(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).Build()
	defer logger.Close()

	logger.Debug().Msg("(API) \"" + r.Pattern + "\"")
}

func ApiCodeStatus(w http.ResponseWriter, statusCode int, message string, errors []rules.ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if message == "" {
		message = httpStatusTexts[statusCode]
	}

	response := map[string]interface{}{
		"status":  httpStatusTexts[statusCode],
		"message": message,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	_ = json.NewEncoder(w).Encode(response)
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
