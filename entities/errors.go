package entities

import (
	"fmt"

	"github.com/powerman/narada-go/narada"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"golang.org/x/net/context"
)

// Service specific API errors
var (
	ErrUnauthorized       = jsonrpc2.NewError(1, "INVALID_TOKEN")
	ErrMaintenance        = jsonrpc2.NewError(4, "MAINTENANCE")
	ErrServerUnavailable  = jsonrpc2.NewError(5, "REMOTE_SERVER_UNAVAILABLE")
	ErrNotFound           = jsonrpc2.NewError(404, "NOT_FOUND")
	ErrInvalidRequest     = jsonrpc2.NewError(101, "TRACKER_VALIDATION_ERROR")
	ErrInvalidCredentials = jsonrpc2.NewError(102, "INVALID_CREDENTIALS")
	ErrInvalidTrackerURL  = jsonrpc2.NewError(104, "INVALID_TRACKER_URL")
	ErrProjectNotFound    = jsonrpc2.NewError(106, "PROJECT_NOT_FOUND")
	ErrIssueNotFound      = jsonrpc2.NewError(107, "ISSUE_NOT_FOUND")
)

// NewServerError creates new JSON RPC error with given message
func NewServerError(msg string) error {
	return jsonrpc2.NewError(-32000, msg)
}

// NewLoggedError logs the error and replaces it with an API error
func NewLoggedError(l *narada.Log, ctx context.Context, err error, msg interface{}) error {
	var test interface{} = err
	_, ok := test.(*jsonrpc2.Error)
	if ok { // error already processed upstream
		return err
	}
	l.ERR(fmt.Sprintf("[%v] %v", ctx.Value("TracingID"), err.Error()))
	if e, ok := msg.(error); ok {
		return e
	}
	if s, ok := msg.(string); ok {
		if s == "" {
			s = err.Error()
		}
		return NewServerError(s)
	}
	return NewServerError("Unexpected error")
}

// NewErrorLogger creates a curries version of NewLoggedError with context and message set
func NewErrorLogger(l *narada.Log, ctx context.Context, msg string) func(error) error {
	return func(err error) error {
		return NewLoggedError(l, ctx, err, msg)
	}
}
