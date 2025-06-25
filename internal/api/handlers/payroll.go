package handlers

import (
	"net/http"
	"payslip/internal/domain/interfaces"

	"github.com/labstack/echo/v4"
)

type PayrollHandler struct {
	payrollService interfaces.PayrollService
}

func NewPayrollHandler(payrollService interfaces.PayrollService) *PayrollHandler {
	return &PayrollHandler{payrollService: payrollService}
}

func (h *PayrollHandler) RunPayroll(c echo.Context) error {
	periodID := c.Param("period_id")

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	if err := h.payrollService.RunPayroll(c.Request().Context(), periodID, userID, c.RealIP(), c.Response().Header().Get(echo.HeaderXRequestID)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Payroll processed"})
}

func (h *PayrollHandler) GeneratePayslip(c echo.Context) error {
	periodID := c.Param("period_id")

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	payslip, err := h.payrollService.GeneratePayslip(c.Request().Context(), periodID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, payslip)
}

func (h *PayrollHandler) GeneratePayrollSummary(c echo.Context) error {
	periodID := c.Param("period_id")

	summary, err := h.payrollService.GeneratePayrollSummary(c.Request().Context(), periodID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, summary)
}
