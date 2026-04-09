package models

import (
	"time"

	"gorm.io/gorm"
)

type Webhook struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	Hash        string            `gorm:"index;not null" json:"hash"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `gorm:"serializer:json;type:jsonb" json:"headers"`
	Body        map[string]string `gorm:"serializer:json;type:jsonb" json:"body"`
	QueryString map[string]string `gorm:"serializer:json;type:jsonb" json:"query_string"`
	ContentType string            `json:"content_type"`
	RemoteAddr  string            `json:"remote_addr"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   time.Time         `gorm:"index" json:"expires_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
}


