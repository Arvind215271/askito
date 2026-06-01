// ./internal/api/errors.go

package api

import (
	"fmt"
	"net/http"
	"time"
)

type Errors struct {
	Common CommonErrors
}

var Err = Errors{
	Common: CommonErrors{},
}

type CommonErrors struct{}

func (CommonErrors) BadRequest(msg string) *AppError {
	return NewError(
		"BAD_REQUEST",
		msg,
		http.StatusBadRequest,
	)
}

func (CommonErrors) Validation() *AppError {
	return NewError(
		"VALIDATION_ERROR",
		"Invalid input",
		http.StatusBadRequest,
	)
}

func (CommonErrors) Unauthorized() *AppError {
	return NewError(
		"UNAUTHORIZED",
		"Unauthorized",
		http.StatusUnauthorized,
	)
}

func (CommonErrors) Forbidden() *AppError {
	return NewError(
		"FORBIDDEN",
		"Forbidden",
		http.StatusForbidden,
	)
}

func (CommonErrors) NotFound(resource string) *AppError {
	return NewError(
		"NOT_FOUND",
		resource+" not found",
		http.StatusNotFound,
	)
}

func (CommonErrors) RouteNotFound() *AppError {
	return NewError(
		"ROUTE_NOT_FOUND",
		"Route not found",
		http.StatusNotFound,
	)
}

func (CommonErrors) Conflict(msg string) *AppError {
	return NewError(
		"CONFLICT",
		msg,
		http.StatusConflict,
	)
}

func (CommonErrors) RateLimited(retryAfter time.Duration) *AppError {
	return NewError(
		"RATE_LIMITED",
		"Too many requests. Please try again later.",
		http.StatusTooManyRequests,
	).AddField(
		"retry_after_seconds",
		fmt.Sprintf("%d", int(retryAfter.Seconds())),
	)
}

func (CommonErrors) Internal() *AppError {
	return NewError(
		"INTERNAL_ERROR",
		"Something went wrong",
		http.StatusInternalServerError,
	)
}