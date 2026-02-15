package models

import "time"

type Address struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Area       string    `json:"area"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	
	// Nullable FKs (Mzyan, hit Address t9der tkoun dyal User aw Bank)
	UserID *uint `json:"user_id"` 
	BankID *uint `json:"bank_id"` 
	
	// ✅ FIX: float64 for GPS Precision
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}