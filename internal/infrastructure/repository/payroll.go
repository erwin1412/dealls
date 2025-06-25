package repository

import (
	"context"
	"fmt"
	"payslip/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PayrollRepository struct {
	db *gorm.DB
}

func NewPayrollRepository(db *gorm.DB) *PayrollRepository {
	return &PayrollRepository{db: db}
}

func (r *PayrollRepository) CreatePayroll(ctx context.Context, payroll *models.Payroll) error {
	return r.db.WithContext(ctx).Create(payroll).Error
}

func (r *PayrollRepository) FindPayrollByPeriodAndUser(ctx context.Context, periodID, userID uuid.UUID) (*models.Payroll, error) {
	var payroll models.Payroll
	if err := r.db.WithContext(ctx).Where("period_id = ? AND user_id = ?", periodID, userID).First(&payroll).Error; err != nil {
		return nil, fmt.Errorf("payroll not found: %w", err)
	}
	return &payroll, nil
}

func (r *PayrollRepository) FindPayrollsByPeriod(ctx context.Context, periodID uuid.UUID) ([]*models.Payroll, error) {
	var payrolls []*models.Payroll
	if err := r.db.WithContext(ctx).Where("period_id = ?", periodID).Find(&payrolls).Error; err != nil {
		return nil, fmt.Errorf("failed to find payrolls: %w", err)
	}
	return payrolls, nil
}

func (r *PayrollRepository) FindAttendancesByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]*models.Attendance, error) {
	var attendances []*models.Attendance
	if err := r.db.WithContext(ctx).Where("user_id = ? AND period_id = ?", userID, periodID).Find(&attendances).Error; err != nil {
		return nil, fmt.Errorf("failed to find attendances: %w", err)
	}
	return attendances, nil
}

func (r *PayrollRepository) FindOvertimesByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]*models.Overtime, error) {
	var overtimes []*models.Overtime
	if err := r.db.WithContext(ctx).Where("user_id = ? AND period_id = ?", userID, periodID).Find(&overtimes).Error; err != nil {
		return nil, fmt.Errorf("failed to find overtimes: %w", err)
	}
	return overtimes, nil
}

func (r *PayrollRepository) FindReimbursementsByUserAndPeriod(ctx context.Context, userID, periodID uuid.UUID) ([]*models.Reimbursement, error) {
	var reimbursements []*models.Reimbursement
	if err := r.db.WithContext(ctx).Where("user_id = ? AND period_id = ?", userID, periodID).Find(&reimbursements).Error; err != nil {
		return nil, fmt.Errorf("failed to find reimbursements: %w", err)
	}
	return reimbursements, nil
}

func (r *PayrollRepository) FindEmployees(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	if err := r.db.WithContext(ctx).Where("role = ?", "employee").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to find employees: %w", err)
	}
	return users, nil
}

func (r *PayrollRepository) CountAttendance(ctx context.Context, userID, periodID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Attendance{}).Where("user_id = ? AND period_id = ?", userID, periodID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count attendance: %w", err)
	}
	return count, nil
}

func (r *PayrollRepository) SumOvertimeHours(ctx context.Context, userID, periodID uuid.UUID) (float64, error) {
	var totalHours float64
	if err := r.db.WithContext(ctx).Model(&models.Overtime{}).Where("user_id = ? AND period_id = ?", userID, periodID).Select("SUM(hours)").Scan(&totalHours).Error; err != nil {
		return 0, fmt.Errorf("failed to sum overtime hours: %w", err)
	}
	return totalHours, nil
}

func (r *PayrollRepository) SumReimbursementAmount(ctx context.Context, userID, periodID uuid.UUID) (float64, error) {
	var totalAmount float64
	if err := r.db.WithContext(ctx).Model(&models.Reimbursement{}).Where("user_id = ? AND period_id = ?", userID, periodID).Select("SUM(amount)").Scan(&totalAmount).Error; err != nil {
		return 0, fmt.Errorf("failed to sum reimbursement amount: %w", err)
	}
	return totalAmount, nil
}

func (r *PayrollRepository) FindUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
