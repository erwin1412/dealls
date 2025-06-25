// internal/infrastructure/repository/audit.go
package repository

import (
	"context"
	"payslip/internal/domain/models"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, audit *models.AuditLog) error {
	return r.db.WithContext(ctx).Create(audit).Error
}
