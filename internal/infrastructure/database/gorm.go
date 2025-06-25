package database

import (
	"payslip/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGORM(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	db.AutoMigrate(
		&models.User{},
		&models.AttendancePeriod{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
		&models.Payroll{},
		&models.AuditLog{},
	)
}
