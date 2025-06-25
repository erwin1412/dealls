package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"payslip/internal/infrastructure/auth"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			log.Printf("RequestID: %s, Method: %s, Path: %s, Status: %d, Duration: %v",
				c.Response().Header().Get(echo.HeaderXRequestID), c.Request().Method, c.Request().URL.Path, c.Response().Status, time.Since(start))
			return err
		}
	}
}

func AuthMiddleware(authService auth.AuthService, requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, role, err := authService.ValidateToken(c.Request().Header.Get("Authorization"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}
			if role != requiredRole {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
			}

			ctx := context.WithValue(c.Request().Context(), userIDKey, userID)
			ctx = context.WithValue(ctx, roleKey, role)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userIDVal := ctx.Value(userIDKey)
	if userIDVal == nil {
		return uuid.UUID{}, fmt.Errorf("user ID not found in context")
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("invalid user ID type in context")
	}
	return userID, nil
}
