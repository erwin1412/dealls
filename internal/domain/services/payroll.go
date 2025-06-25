package services

import (
	"context"
	"fmt"
	"payslip/internal/domain/interfaces"
	"payslip/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type PayrollService struct {
	payrollRepo    interfaces.PayrollRepository
	attendanceRepo interfaces.AttendanceRepository
	auditRepo      interfaces.AuditRepository
}

func NewPayrollService(payrollRepo interfaces.PayrollRepository, attendanceRepo interfaces.AttendanceRepository, auditRepo interfaces.AuditRepository) *PayrollService {
	return &PayrollService{payrollRepo: payrollRepo, attendanceRepo: attendanceRepo, auditRepo: auditRepo}
}

func (s *PayrollService) RunPayroll(ctx context.Context, periodID string, userID uuid.UUID, ipAddress, requestID string) error {
	parsedPeriodID, err := uuid.Parse(periodID)
	if err != nil {
		return fmt.Errorf("invalid period ID: %w", err)
	}

	period, err := s.attendanceRepo.FindPeriodByID(ctx, parsedPeriodID)
	if err != nil {
		return fmt.Errorf("period not found: %w", err)
	}

	if s.attendanceRepo.IsPayrollProcessed(ctx, parsedPeriodID) {
		return fmt.Errorf("payroll already processed for this period")
	}

	workingDays := countWorkingDays(period.StartDate, period.EndDate)
	totalWorkingHours := float64(workingDays * 8)

	employees, err := s.payrollRepo.FindEmployees(ctx)
	if err != nil {
		return fmt.Errorf("failed to find employees: %w", err)
	}

	for _, user := range employees {
		attendanceCount, err := s.payrollRepo.CountAttendance(ctx, user.ID, parsedPeriodID)
		if err != nil {
			return fmt.Errorf("failed to count attendance for user %s: %w", user.ID, err)
		}

		salaryPerHour := user.Salary / totalWorkingHours
		baseSalary := salaryPerHour * float64(attendanceCount*8)

		totalOvertimeHours, err := s.payrollRepo.SumOvertimeHours(ctx, user.ID, parsedPeriodID)
		if err != nil {
			return fmt.Errorf("failed to sum overtime for user %s: %w", user.ID, err)
		}
		overtimePay := salaryPerHour * 2 * totalOvertimeHours

		totalReimbursement, err := s.payrollRepo.SumReimbursementAmount(ctx, user.ID, parsedPeriodID)
		if err != nil {
			return fmt.Errorf("failed to sum reimbursement for user %s: %w", user.ID, err)
		}

		totalPay := baseSalary + overtimePay + totalReimbursement

		payroll := &models.Payroll{
			ID:                  uuid.New(),
			PeriodID:            parsedPeriodID,
			UserID:              user.ID,
			BaseSalary:          baseSalary,
			OvertimePay:         overtimePay,
			ReimbursementAmount: totalReimbursement,
			TotalPay:            totalPay,
			CreatedBy:           userID,
			IPAddress:           ipAddress,
		}

		if err := s.payrollRepo.CreatePayroll(ctx, payroll); err != nil {
			return fmt.Errorf("failed to create payroll for user %s: %w", user.ID, err)
		}

		audit := &models.AuditLog{
			ID:        uuid.New(),
			Action:    "create",
			TableName: "payroll",
			RecordID:  payroll.ID,
			UserID:    userID,
			IPAddress: ipAddress,
			RequestID: requestID,
			Details:   fmt.Sprintf("Processed payroll for user %s for period %s", user.ID, periodID),
			CreatedAt: time.Now(),
		}
		if err := s.auditRepo.Create(ctx, audit); err != nil {
			return fmt.Errorf("failed to log audit: %w", err)
		}
	}

	return nil
}

func (s *PayrollService) GeneratePayslip(ctx context.Context, periodID string, userID uuid.UUID) (map[string]interface{}, error) {
	parsedPeriodID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	payroll, err := s.payrollRepo.FindPayrollByPeriodAndUser(ctx, parsedPeriodID, userID)
	if err != nil {
		return nil, fmt.Errorf("payroll not found: %w", err)
	}

	period, err := s.attendanceRepo.FindPeriodByID(ctx, parsedPeriodID)
	if err != nil {
		return nil, fmt.Errorf("period not found: %w", err)
	}

	attendances, err := s.payrollRepo.FindAttendancesByUserAndPeriod(ctx, userID, parsedPeriodID)
	if err != nil {
		return nil, fmt.Errorf("failed to find attendances: %w", err)
	}

	overtimes, err := s.payrollRepo.FindOvertimesByUserAndPeriod(ctx, userID, parsedPeriodID)
	if err != nil {
		return nil, fmt.Errorf("failed to find overtimes: %w", err)
	}

	reimbursements, err := s.payrollRepo.FindReimbursementsByUserAndPeriod(ctx, userID, parsedPeriodID)
	if err != nil {
		return nil, fmt.Errorf("failed to find reimbursements: %w", err)
	}

	return map[string]interface{}{
		"period":               period,
		"attendance":           attendances,
		"overtime":             overtimes,
		"reimbursements":       reimbursements,
		"base_salary":          payroll.BaseSalary,
		"overtime_pay":         payroll.OvertimePay,
		"reimbursement_amount": payroll.ReimbursementAmount,
		"total_pay":            payroll.TotalPay,
	}, nil
}

func (s *PayrollService) GeneratePayrollSummary(ctx context.Context, periodID string) (map[string]interface{}, error) {
	parsedPeriodID, err := uuid.Parse(periodID)
	if err != nil {
		return nil, fmt.Errorf("invalid period ID: %w", err)
	}

	payrolls, err := s.payrollRepo.FindPayrollsByPeriod(ctx, parsedPeriodID)
	if err != nil {
		return nil, fmt.Errorf("failed to find payrolls: %w", err)
	}

	var totalPay float64
	summary := make([]map[string]interface{}, len(payrolls))
	for i, p := range payrolls {
		user, err := s.payrollRepo.FindUserByID(ctx, p.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user %s: %w", p.UserID, err)
		}
		summary[i] = map[string]interface{}{
			"username":  user.Username,
			"total_pay": p.TotalPay,
		}
		totalPay += p.TotalPay
	}

	return map[string]interface{}{
		"summary":       summary,
		"total_payroll": totalPay,
	}, nil
}

func countWorkingDays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			count++
		}
	}
	return count
}
