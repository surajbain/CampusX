package models

import (
	"time"

	"gorm.io/gorm"
)

// Base is embedded in every model that supports soft delete.
// GORM will look for a `deleted_at` column — migrations MUST create it.
type Base struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BaseNoSoftDelete for tables without soft delete.
type BaseNoSoftDelete struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
