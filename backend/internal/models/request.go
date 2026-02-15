package models

import "time"

// Status Enum
type RequestStatus string

const (
	StatusPending   RequestStatus = "Pending"   // Yallah t-hat (Tsnna validation)
	StatusApproved  RequestStatus = "Approved"  // Admin qblo (Notification mchat)
	StatusFulfilled RequestStatus = "Fulfilled" // Chi wa7ed tbrr3
	StatusCancelled RequestStatus = "Cancelled"
)

type BloodRequest struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Chkon Talab?
	UserID uint `json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`

	// Achmn Dam?
	BloodType string `json:"blood_type"` // e.g. "O+"
	IsUrgent  bool   `json:"is_urgent"`  // Wach 7ala khatira?

	// Fin? (Location dyal Sbitar)
	HospitalName string  `json:"hospital_name"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`

	// Daleel (Proof)
	PrescriptionPhoto string `json:"prescription_photo"` // URL dyal tswira

	// 7ala
	Status RequestStatus `json:"status" gorm:"default:'Pending'"`

	CreatedAt time.Time `json:"created_at"`
}