package interfaces

import (
	"context"
	"payslip/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type AttendanceRepository interface {
	CreatePeriod(ctx context.Context, period *models.AttendancePeriod) error
	FindPeriodByID(ctx context.Context, id uuid.UUID) (*models.AttendancePeriod, error)
	CreateAttendance(ctx context.Context, attendance *models.Attendance) error
	FindAttendanceByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time, periodID uuid.UUID) (*models.Attendance, error)
	CreateOvertime(ctx context.Context, overtime *models.Overtime) error
	CreateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) error
	IsPayrollProcessed(ctx context.Context, periodID uuid.UUID) bool
}

type AttendanceService interface {
	CreatePeriod(ctx context.Context, startDate, endDate string, userID uuid.UUID, ipAddress, requestID string) (*models.AttendancePeriod, error)
	SubmitAttendance(ctx context.Context, date, periodID string, userID uuid.UUID, ipAddress, requestID string) (*models.Attendance, error)
	SubmitOvertime(ctx context.Context, date string, hours float64, periodID string, userID uuid.UUID, ipAddress, requestID string) (*models.Overtime, error)
	SubmitReimbursement(ctx context.Context, amount float64, description, periodID string, userID uuid.UUID, ipAddress, requestID string) (*models.Reimbursement, error)
}
