// internal/domain/models/user.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"unique;not null;size:50"`
	Password  string    `gorm:"not null;size:100"`
	Role      string    `gorm:"not null;size:20"` // 'employee' or 'admin'
	Salary    float64   `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID
	UpdatedBy uuid.UUID
}
