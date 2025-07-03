package model

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Status      string         `json:"status" gorm:"default:doing,not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}

func (Todo) TableName() string {
	return "todo"
}
