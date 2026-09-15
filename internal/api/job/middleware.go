package job

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

// UserIdentityMiddleware validates the X-User-ID header as a valid UUID and stores it in the request context.
func UserIdentityMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userIDStr := (*c).Request().Header.Get("X-User-ID")
			if userIDStr == "" {
				return Err.MissingOrInvalidUserID()
			}

			parsed, err := uuid.Parse(userIDStr)
			if err != nil {
				return Err.MissingOrInvalidUserID()
			}

			ctx := context.WithValue((*c).Request().Context(), userIDKey, parsed.String())
			(*c).SetRequest((*c).Request().WithContext(ctx))

			return next(c)
		}
	}
}

// UserIDFromContext retrieves the validated user ID from the context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
