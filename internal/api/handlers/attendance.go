package handlers

import (
	"net/http"
	"payslip/internal/domain/interfaces"

	"github.com/labstack/echo/v4"
)

type AttendanceHandler struct {
	attendanceService interfaces.AttendanceService
}

func NewAttendanceHandler(attendanceService interfaces.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: attendanceService}
}

func (h *AttendanceHandler) CreateAttendancePeriod(c echo.Context) error {
	var input struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	period, err := h.attendanceService.CreatePeriod(c.Request().Context(), input.StartDate, input.EndDate, userID, c.RealIP(), c.Response().Header().Get(echo.HeaderXRequestID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Period created",
		"period_id": period.ID,
	})
}

func (h *AttendanceHandler) SubmitAttendance(c echo.Context) error {
	var input struct {
		Date     string `json:"date"`
		PeriodID string `json:"period_id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	attendance, err := h.attendanceService.SubmitAttendance(c.Request().Context(), input.Date, input.PeriodID, userID, c.RealIP(), c.Response().Header().Get(echo.HeaderXRequestID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Attendance submitted",
		"attendance_id": attendance.ID,
	})
}

func (h *AttendanceHandler) SubmitOvertime(c echo.Context) error {
	var input struct {
		Date     string  `json:"date"`
		Hours    float64 `json:"hours"`
		PeriodID string  `json:"period_id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	overtime, err := h.attendanceService.SubmitOvertime(c.Request().Context(), input.Date, input.Hours, input.PeriodID, userID, c.RealIP(), c.Response().Header().Get(echo.HeaderXRequestID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Overtime submitted",
		"overtime_id": overtime.ID,
	})
}

func (h *AttendanceHandler) SubmitReimbursementByID(c echo.Context) error {
	var input struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		PeriodID    string  `json:"period_id"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userID, err := GetUserIDFromContext(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	reimbursement, err := h.attendanceService.SubmitReimbursement(c.Request().Context(), input.Amount, input.Description, input.PeriodID, userID, c.RealIP(), c.Response().Header().Get(echo.HeaderXRequestID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":          "Reimbursement submitted",
		"reimbursement_id": reimbursement.ID,
	})
}
