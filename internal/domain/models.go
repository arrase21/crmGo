package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	FirstName string         `gorm:"size:30;not null" json:"first_name"`
	LastName  string         `gorm:"size:40;not null" json:"last_name"`
	Dni       string         `gorm:"size:20;not null;uniqueIndex" json:"dni"`
	Gender    string         `gorm:"size:3;not null" json:"gender"`
	Phone     string         `gorm:"size:15;not null;uniqueIndex" json:"phone"`
	Email     string         `gorm:"size:50;not null;uniqueIndex" json:"email"`
	BirthDay  time.Time      `gorm:"not null" json:"birth_day"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// type Permission struct {
// 	ID          uint      `gorm:"primaryKey" json:"id"`
// 	Name        string    `gorm:"size:40;not null" json:"name"`
// 	Description string    `gorm:"size:50;not null" json:"description"`
// 	Actions     []string  `gorm:"size:"`
// 	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
// 	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
// 	DeletedAt   time.Time `gorm:"not null" json:"deleted_at"`
// }

// type PermissionAction struct {
// 	ID uint `gorm:"primaryKey" json:"id"`
// 	PermissionID
// }
