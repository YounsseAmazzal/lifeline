package models

import "time"

type Address struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Area       string    `json:"area"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	
	UserID *uint `json:"user_id"` // Nullable
	BankID *uint `json:"bank_id"` // Nullable
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Longitude float32 `json:"longitude"`
	Latitude float32 `json:"latitude"`
}