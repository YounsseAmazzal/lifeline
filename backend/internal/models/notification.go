package models

import "time"

type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"` // Lamin msiftha?
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`    // "Urgent", "Info", "Success"
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}