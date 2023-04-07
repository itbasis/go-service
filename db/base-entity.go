package db

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type BaseEntity struct {
	UID *uuid.UUID `gorm:"column:uid; type:uuid; primaryKey; default:uuid_generate_v4()"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
