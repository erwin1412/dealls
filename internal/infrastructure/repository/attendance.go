package repository

import (
	"context"
	"fmt"
	"payslip/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) CreatePeriod(ctx context.Context, period *models.AttendancePeriod) error {
	return r.db.WithContext(ctx).Create(period).Error
}

func (r *AttendanceRepository) FindPeriodByID(ctx context.Context, id uuid.UUID) (*models.AttendancePeriod, error) {
	var period models.AttendancePeriod
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&period).Error; err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}
	return &period, nil
}

func (r *AttendanceRepository) CreateAttendance(ctx context.Context, attendance *models.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

func (r *AttendanceRepository) FindAttendanceByUserAndDate(ctx context.Context, userID uuid.UUID, date time.Time, periodID uuid.UUID) (*models.Attendance, error) {
	var attendance models.Attendance
	if err := r.db.WithContext(ctx).Where("user_id = ? AND date = ? AND period_id = ?", userID, date, periodID).First(&attendance).Error; err != nil {
		return nil, fmt.Errorf("attendance not found: %w", err)
	}
	return &attendance, nil
}

func (r *AttendanceRepository) CreateOvertime(ctx context.Context, overtime *models.Overtime) error {
	return r.db.WithContext(ctx).Create(overtime).Error
}

func (r *AttendanceRepository) CreateReimbursement(ctx context.Context, reimbursement *models.Reimbursement) error {
	return r.db.WithContext(ctx).Create(reimbursement).Error
}

func (r *AttendanceRepository) IsPayrollProcessed(ctx context.Context, periodID uuid.UUID) bool {
	var count int64
	r.db.WithContext(ctx).Model(&models.Payroll{}).Where("period_id = ?", periodID).Count(&count)
	return count > 0
}
