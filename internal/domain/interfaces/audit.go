// internal/domain/interfaces/audit.go
package interfaces

import (
	"context"
	"payslip/internal/domain/models"
)

type AuditRepository interface {
	Create(ctx context.Context, audit *models.AuditLog) error
}
