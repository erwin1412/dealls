package models

import (
	"time"

	"github.com/google/uuid"
)

type Payroll struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PeriodID            uuid.UUID `gorm:"not null"`
	UserID              uuid.UUID `gorm:"not null"`
	BaseSalary          float64   `gorm:"not null"`
	OvertimePay         float64   `gorm:"not null"`
	ReimbursementAmount float64   `gorm:"not null"`
	TotalPay            float64   `gorm:"not null"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	CreatedBy           uuid.UUID
	IPAddress           string `gorm:"size:45"`
}
