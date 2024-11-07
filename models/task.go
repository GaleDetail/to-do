package models

import "time"

type Task struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Status      string    `gorm:"type:varchar(1000)"`
	Description string    `gorm:"type:enum('pending','in_progress','completed');not null; default:'pending'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
