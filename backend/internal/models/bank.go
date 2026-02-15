package models

import (
	"time"
)

type Bank struct {
	ID          uint      `json:"id"`
	City     string      `json:"city"`
	Name        string    `json:"name"` // Correction: name (sghira)
	PhoneNumber string    `json:"phone_number"`         // Correction: phone_number
	Email       string    `json:"email"`  // Correction: email
	Website     string    `json:"website"`              // Correction: website
	LastUpdated time.Time `json:"last_updated"`

	// --- RELATIONSHIPS ---

	Address Address `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address"`

	Photo Photo `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo"`

	// Stock
	BloodGroups []BloodGroup `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"blood_groups"`

	// Staff
	Moderators []Moderator `gorm:"foreignKey:BankID" json:"moderators"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Constructor (Helper Function)
func NewBank(city,name, phone, email string) Bank {
	bloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	
	var initialGroups []BloodGroup
	for _, bt := range bloodTypes {
		initialGroups = append(initialGroups, BloodGroup{
			Group:    bt,
			Quantity: 0,
		})
	}

	return Bank{
		City: 			city,
		Name:        name,
		PhoneNumber: phone,
		Email:       email,
		BloodGroups: initialGroups,
		LastUpdated: time.Now(),
	}
}