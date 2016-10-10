package entities

import (
	"net/http"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

//Global error codes
var (
	ErrTimeout      = jsonrpc2.NewError(0, "TIMEOUT")
	ErrForbidden    = jsonrpc2.NewError(3, "FORBIDDEN")
	ErrMaintenance  = jsonrpc2.NewError(4, "MAINTENANCE")
	ErrRemoteServer = jsonrpc2.NewError(5, "REMOTE_SERVER_UNAVAILABLE")
)

//Tracker services error codes
var (
	ErrCredentials     = jsonrpc2.NewError(102, "INVALID_CREDENTIALS")
	ErrTrackerType     = jsonrpc2.NewError(103, "INVALID_TRACKER_TYPE")
	ErrTrackerURL      = jsonrpc2.NewError(104, "INVALID_TRACKER_URL")
	ErrIssueURL        = jsonrpc2.NewError(105, "INVALID_ISSUE_URL")
	ErrProjectNotFound = jsonrpc2.NewError(106, "PROJECT_NOT_FOUND")
	ErrIssueNotFound   = jsonrpc2.NewError(107, "ISSUE_NOT_FOUND")
)

const (
	trackerValidationErrCode = 101
	trackerValidationErrMsg  = "TRACKER_VALIDATION_ERROR"
)

//NewTrackerValidationErr return new tracker validation error with message
func NewTrackerValidationErr(msg string) error {
	return &jsonrpc2.Error{
		Code:    trackerValidationErrCode,
		Message: trackerValidationErrMsg,
		Data:    msg,
	}
}

// ErrorForStatusCode convert http.Response status code to our errors
func ErrorForStatusCode(code int) error {
	switch code {
	case http.StatusRequestTimeout:
		return ErrTimeout
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusUnauthorized:
		return ErrCredentials
	case http.StatusServiceUnavailable:
		return ErrRemoteServer
	case http.StatusInternalServerError:
		return ErrRemoteServer
	default:
		return nil
	}
}
