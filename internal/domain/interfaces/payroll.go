package interfaces

import (
	"context"
	"payslip/internal/domain/models"

	"github.com/google/uuid"
)

type PayrollRepository interface {
	CreatePayroll(ctx context.Context, payroll *models.Payroll) error
	FindPayrollByPeriodAndUser(ctx context.Context, periodID, userID uuid.UUID) (*models.Payroll, error)
	FindPayrollsByPeriod(ctx context.Context, periodID uuid.UUID) ([]*models.Payroll, error)
	FindAttendancesByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]*models.Attendance, error)
	FindOvertimesByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]*models.Overtime, error)
	FindReimbursementsByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]*models.Reimbursement, error)
	FindEmployees(ctx context.Context) ([]*models.User, error)
	CountAttendance(ctx context.Context, userID, periodID uuid.UUID) (int64, error)
	SumOvertimeHours(ctx context.Context, userID, periodID uuid.UUID) (float64, error)
	SumReimbursementAmount(ctx context.Context, userID, periodID uuid.UUID) (float64, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) // Added
}
type PayrollService interface {
	RunPayroll(ctx context.Context, periodID string, userID uuid.UUID, ipAddress, requestID string) error
	GeneratePayslip(ctx context.Context, periodID string, userID uuid.UUID) (map[string]interface{}, error)
	GeneratePayrollSummary(ctx context.Context, periodID string) (map[string]interface{}, error)
}
