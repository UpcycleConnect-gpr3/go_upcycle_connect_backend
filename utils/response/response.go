package response

import (
	"encoding/json"
	"go-upcycle_connect-backend/utils/rules"
	"net/http"
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
	ErrAuthFailed           = "Credentials do not match"
	ErrInvalidBody          = "Incorrect body format"
	ErrJson                 = "Json parse error"
	ErrGenerateToken        = "Generate token error"
	ErrUserNotFound         = "User not found"
	ErrInvalidTOTP          = "Invalid TOTP code"
	ErrEnableTOTP           = "Failed to enable TOTP"
	ErrGenerateTOTP         = "Failed to generate TOTP"
	ErrDisableTOTP          = "Failed to disable TOTP"
	ErrAuthTokenRequired    = "Authorization token required"
	ErrInvalidAuthToken     = "Invalid authorization token"
	ErrInvalidOrExpiredHash = "Invalid or expired hash"
	ErrFetchingTOTPRecord   = "Failed to fetch TOTP record"
)

const (
	SuccessEnableTOTP   = "TOTP enabled successfully"
	SuccessDisableTOTP  = "TOTP disabled successfully"
	SuccessGenerateTOTP = "TOTP key generated successfully"
	SuccessLogin        = "Login successful"
)

type Response struct {
	Success          bool                     `json:"success"`
	Data             interface{}              `json:"data,omitempty"`
	Message          string                   `json:"message,omitempty"`
	Errors           *[]rules.ValidationError `json:"errors,omitempty"`
	Error            *ErrorResponse           `json:"error,omitempty"`
	SuccessResponses *[]SuccessResponse       `json:"successResponses,omitempty"`
	SuccessResponse  *SuccessResponse         `json:"successResponse,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func newSuccess(w http.ResponseWriter, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func NewSuccessMessage(w http.ResponseWriter, message string) {
	newSuccess(w, nil, message)
}

func NewSuccessData(w http.ResponseWriter, data interface{}) {
	newSuccess(w, data, "")
}

func newError(w http.ResponseWriter, statusCode int, message string, errors []rules.ValidationError) {
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

func NewValidationError(w http.ResponseWriter, message string, errors []rules.ValidationError) {
	newError(w, http.StatusBadRequest, message, errors)
}

func NewErrorMessage(w http.ResponseWriter, message string, statusCode int) {
	newError(w, statusCode, message, nil)
}

func (r Response) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	if !r.Success {
		status = http.StatusBadRequest
		if r.Error != nil {
			status = http.StatusInternalServerError
		}
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(r)
}

func NewMessage(message string) Response {
	return Response{
		Success: true,
		Message: message,
	}
}
