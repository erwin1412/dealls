package models

import (
	"time"

	"github.com/google/uuid"
)

type AttendancePeriod struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	StartDate time.Time `gorm:"not null;type:date"`
	EndDate   time.Time `gorm:"not null;type:date"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}

type Attendance struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	Date      time.Time `gorm:"not null;type:date"`
	PeriodID  uuid.UUID `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	IPAddress string `gorm:"size:45"`
}

type Overtime struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	Date      time.Time `gorm:"not null;type:date"`
	Hours     float64   `gorm:"not null"`
	PeriodID  uuid.UUID `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
	IPAddress string `gorm:"size:45"`
}

type Reimbursement struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID      uuid.UUID `gorm:"not null"`
	Amount      float64   `gorm:"not null"`
	Description string    `gorm:"not null;type:text"`
	PeriodID    uuid.UUID `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
	IPAddress   string `gorm:"size:45"`
}
