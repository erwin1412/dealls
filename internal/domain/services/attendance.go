package services

import (
	"context"
	"fmt"
	"payslip/internal/domain/interfaces"
	"payslip/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type AttendanceService struct {
	attendanceRepo interfaces.AttendanceRepository
	auditRepo      interfaces.AuditRepository
}

func NewAttendanceService(attendanceRepo interfaces.AttendanceRepository, auditRepo interfaces.AuditRepository) *AttendanceService {
	return &AttendanceService{attendanceRepo: attendanceRepo, auditRepo: auditRepo}
}

func (s *AttendanceService) CreatePeriod(ctx context.Context, startDate, endDate string, userID uuid.UUID, ipAddress, requestID string) (*models.AttendancePeriod, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	period := &models.AttendancePeriod{
		ID:        uuid.New(),
		StartDate: start,
		EndDate:   end,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	err = s.attendanceRepo.CreatePeriod(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("failed to create period: %w", err)
	}

	audit := &models.AuditLog{
		ID:        uuid.New(),
		Action:    "create",
		TableName: "attendance_period",
		RecordID:  period.ID,
		UserID:    userID,
		IPAddress: ipAddress,
		RequestID: requestID,
		Details:   fmt.Sprintf("Created attendance period %s from %s to %s", period.ID, startDate, endDate),
		CreatedAt: time.Now(),
	}
	if err := s.auditRepo.Create(ctx, audit); err != nil {
		return nil, fmt.Errorf("failed to log audit: %w", err)
	}

	return period, nil
}

func (s *AttendanceService) SubmitAttendance(ctx context.Context, date, periodID string, userID uuid.UUID, ipAddress, requestID string) (*models.Attendance, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	parsedPeriodID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	if s.attendanceRepo.IsPayrollProcessed(ctx, parsedPeriodID) {
		return nil, fmt.Errorf("payroll already processed for this period")
	}

	if parsedDate.Weekday() == time.Saturday || parsedDate.Weekday() == time.Sunday {
		return nil, fmt.Errorf("cannot submit attendance on weekends")
	}

	if _, err := s.attendanceRepo.FindAttendanceByUserAndDate(ctx, userID, parsedDate, parsedPeriodID); err == nil {
		return nil, fmt.Errorf("attendance already submitted for this date")
	}

	attendance := &models.Attendance{
		ID:        uuid.New(),
		UserID:    userID,
		Date:      parsedDate,
		PeriodID:  parsedPeriodID,
		CreatedBy: userID,
		UpdatedBy: userID,
		IPAddress: ipAddress,
	}

	if err := s.attendanceRepo.CreateAttendance(ctx, attendance); err != nil {
		return nil, fmt.Errorf("failed to submit attendance: %v", err)
	}

	audit := &models.AuditLog{
		ID:        uuid.New(),
		Action:    "create",
		TableName: "attendance",
		RecordID:  attendance.ID,
		UserID:    userID,
		IPAddress: ipAddress,
		RequestID: requestID,
		Details:   fmt.Sprintf("Submitted attendance for user %s on %s", userID, date),
		CreatedAt: time.Now(),
	}
	if err := s.auditRepo.Create(ctx, audit); err != nil {
		return nil, fmt.Errorf("failed to log audit: %v", err)
	}

	return attendance, nil
}

func (s *AttendanceService) SubmitOvertime(ctx context.Context, date string, hours float64, periodID string, userID uuid.UUID, ipAddress, requestID string) (*models.Overtime, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	parsedPeriodID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	if s.attendanceRepo.IsPayrollProcessed(ctx, parsedPeriodID) {
		return nil, fmt.Errorf("payroll already processed for this period")
	}

	if hours > 3 {
		return nil, fmt.Errorf("overtime cannot exceed 3 hours per day")
	}

	overtime := &models.Overtime{
		ID:        uuid.New(),
		UserID:    userID,
		Date:      parsedDate,
		Hours:     hours,
		PeriodID:  parsedPeriodID,
		CreatedBy: userID,
		UpdatedBy: userID,
		IPAddress: ipAddress,
	}

	if err := s.attendanceRepo.CreateOvertime(ctx, overtime); err != nil {
		return nil, fmt.Errorf("failed to submit overtime: %v", err)
	}

	audit := &models.AuditLog{
		ID:        uuid.New(),
		Action:    "create",
		TableName: "overtime",
		RecordID:  overtime.ID,
		UserID:    userID,
		IPAddress: ipAddress,
		RequestID: requestID,
		Details:   fmt.Sprintf("Submitted %f hours overtime for user %s on %s", hours, userID, date),
		CreatedAt: time.Now(),
	}
	if err := s.auditRepo.Create(ctx, audit); err != nil {
		return nil, fmt.Errorf("failed to log audit: %v", err)
	}

	return overtime, nil
}

func (s *AttendanceService) SubmitReimbursement(ctx context.Context, amount float64, description, periodID string, userID uuid.UUID, ipAddress, requestID string) (*models.Reimbursement, error) {
	parsedPeriodID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	if s.attendanceRepo.IsPayrollProcessed(ctx, parsedPeriodID) {
		return nil, fmt.Errorf("payroll already processed for this period")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}
	if description == "" {
		return nil, fmt.Errorf("description is required")
	}

	reimbursement := &models.Reimbursement{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      amount,
		Description: description,
		PeriodID:    parsedPeriodID,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		IPAddress:   ipAddress,
	}

	if err := s.attendanceRepo.CreateReimbursement(ctx, reimbursement); err != nil {
		return nil, fmt.Errorf("failed to submit reimbursement: %v", err)
	}

	audit := &models.AuditLog{
		ID:        uuid.New(),
		Action:    "create",
		TableName: "reimbursement",
		RecordID:  reimbursement.ID,
		UserID:    userID,
		IPAddress: ipAddress,
		RequestID: requestID,
		Details:   fmt.Sprintf("Submitted reimbursement of $%f for user %s", amount, userID),
		CreatedAt: time.Now(),
	}
	if err := s.auditRepo.Create(ctx, audit); err != nil {
		return nil, fmt.Errorf("failed to log audit: %v", err)
	}

	return reimbursement, nil
}
