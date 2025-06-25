// internal/domain/interfaces/user.go
package interfaces

import (
	"context"
	"payslip/internal/domain/models"
)

type UserService interface {
	Register(ctx context.Context, username, password, role, adminIDStr, ipAddress, requestID string) (*models.User, error)
	Login(ctx context.Context, username, password string) (*models.User, string, error)
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	WithTransaction(ctx context.Context, fn func(tx context.Context) error) error
}
