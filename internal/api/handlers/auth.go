// internal/api/handlers/auth.go
package handlers

import (
	"net/http"
	"payslip/internal/domain/interfaces"
	"payslip/internal/infrastructure/auth"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userService interfaces.UserService
	authService auth.AuthService
}

func NewAuthHandler(userService interfaces.UserService, authService auth.AuthService) *AuthHandler {
	return &AuthHandler{userService: userService, authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	user, err := h.userService.Register(c.Request().Context(), input.Username, input.Password, input.Role, userID.String(), c.RealIP(), c.Response().Header().Get(echo.HeaderXRequestID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "User registered successfully",
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	user, _, err := h.userService.Login(c.Request().Context(), input.Username, input.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	token, err := h.authService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user_id": user.ID,
		"role":    user.Role,
	})
}
