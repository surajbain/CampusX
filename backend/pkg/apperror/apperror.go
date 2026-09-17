package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type Code string

const (
	CodeBadRequest    Code = "BAD_REQUEST"
	CodeUnauthorized  Code = "UNAUTHORIZED"
	CodeForbidden     Code = "FORBIDDEN"
	CodeNotFound      Code = "NOT_FOUND"
	CodeConflict      Code = "CONFLICT"
	CodeUnprocessable Code = "UNPROCESSABLE"
	CodeTooManyReqs   Code = "TOO_MANY_REQUESTS"
	CodeInternal      Code = "INTERNAL_ERROR"
	CodeUnavailable   Code = "SERVICE_UNAVAILABLE"
)

type AppError struct {
	Code    Code
	Message string
	Details map[string]any
	Err     error // wrapped internal error (never exposed to client)
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case CodeBadRequest, CodeUnprocessable:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeTooManyReqs:
		return http.StatusTooManyRequests
	case CodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// ---- constructors ----

func New(code Code, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

func WithDetails(code Code, msg string, d map[string]any) *AppError {
	return &AppError{Code: code, Message: msg, Details: d}
}

func Wrap(code Code, msg string, err error) *AppError {
	return &AppError{Code: code, Message: msg, Err: err}
}

func BadRequest(msg string) *AppError          { return New(CodeBadRequest, msg) }
func Unauthorized(msg string) *AppError        { return New(CodeUnauthorized, msg) }
func Forbidden(msg string) *AppError           { return New(CodeForbidden, msg) }
func NotFound(msg string) *AppError            { return New(CodeNotFound, msg) }
func Conflict(msg string) *AppError            { return New(CodeConflict, msg) }
func Internal(msg string, err error) *AppError { return Wrap(CodeInternal, msg, err) }

// As extracts an *AppError from any error chain, or wraps unknown errors.
func As(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return Internal("internal error", err)
}
