// ./internal/api/response.go
package api

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

// generic response
func JSON(
	c *echo.Context,
	status int,
	message string,
	data any,
	meta any,
) error {

	return c.JSON(status, Response{
		Success: status < 400,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// success response
func Success(
	c *echo.Context,
	message string,
	data any,
	meta any,
) error {

	return JSON(
		c,
		http.StatusOK,
		message,
		data,
		meta,
	)
}

// create error response body
func ErrorResponse(appErr *AppError) Response {

	return Response{
		Success: false,
		Message: appErr.Message,
		Data:    appErr.Fields,
		Meta: map[string]any{
			"code": appErr.Code,
		},
	}
}