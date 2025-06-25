package services

import (
	"context"
	"fmt"
	"payslip/internal/domain/interfaces"
	"payslip/internal/domain/models"
	"regexp"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo  interfaces.UserRepository
	auditRepo interfaces.AuditRepository
}

func NewUserService(userRepo interfaces.UserRepository, auditRepo interfaces.AuditRepository) *UserService {
	return &UserService{userRepo: userRepo, auditRepo: auditRepo}
}

func (s *UserService) Register(ctx context.Context, username, password, role, adminIDStr, ipAddress, requestID string) (*models.User, error) {
	// Validate input
	username = strings.TrimSpace(username)
	role = strings.ToLower(role)
	if username == "" || password == "" || (role != "employee" && role != "admin") {
		return nil, fmt.Errorf("username, password, and valid role (employee or admin) are required")
	}
	if len(password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters long")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return nil, fmt.Errorf("username must be alphanumeric")
	}

	// Check for duplicate
	if _, err := s.userRepo.FindByUsername(ctx, username); err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Parse admin ID
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid admin ID")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:        uuid.New(),
		Username:  username,
		Password:  string(hash),
		Role:      role,
		CreatedBy: adminID,
		UpdatedBy: adminID,
	}
	if role == "employee" {
		user.Salary = gofakeit.Float64Range(2000, 10000)
	}

	// Save user and audit log with transaction
	err = s.userRepo.WithTransaction(ctx, func(tx context.Context) error {
		if err := s.userRepo.Create(tx, user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		audit := &models.AuditLog{
			ID:        uuid.New(),
			Action:    "create",
			TableName: "user",
			RecordID:  user.ID,
			UserID:    adminID,
			IPAddress: ipAddress,
			RequestID: requestID,
			Details:   fmt.Sprintf("Registered user %s with role %s", username, role),
		}
		return s.auditRepo.Create(tx, audit)
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	return user, "", nil // Token generation moved to auth package
}
