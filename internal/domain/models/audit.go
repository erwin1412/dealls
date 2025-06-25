package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Action    string    `gorm:"not null;size:100"`
	TableName string    `gorm:"not null;size:50"`
	RecordID  uuid.UUID
	UserID    uuid.UUID
	IPAddress string    `gorm:"size:45"`
	RequestID string    `gorm:"size:36"`
	Details   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
