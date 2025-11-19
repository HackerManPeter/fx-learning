package models

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;not null" json:"id"`
	CreatedAt time.Time      `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;not null" json:"updated_at"`
	DeletedAt *time.Time     `gorm:"type:timestamp" json:"deleted_at,omitempty"`
	MetaData  map[string]any `gorm:"serializer:json;type:json" json:"meta_data,omitempty"`
}
