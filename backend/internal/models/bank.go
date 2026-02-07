package models

import (
	"time"
)

type Bank struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `gorm:"unique" json:"email"`
	Website     string    `json:"website"`
	LastUpdated time.Time `json:"last_updated"`

	// --- RELATIONSHIPS ---

	Address Address `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"address"`

	Photo Photo `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo"`

	//  Stock 
	BloodGroups []BloodGroup `gorm:"foreignKey:BankID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"blood_groups"`

	// 4. Staff (One-to-Many)
	// Moderators homa l-mwdafin li kiy-jeriw l-bank
	Moderators []Moderator `gorm:"foreignKey:BankID" json:"moderators"`

	// 5. Requests (One-to-Many)
	// Ila bghiti t-3ref chkon talab d-dam mn had l-bank
	// Requests []Request `gorm:"foreignKey:BankID" json:"requests"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Constructor (Helper Function)
func NewBank(name, phone, email string) Bank {
	// Standard 8 Blood Types
	bloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	
	var initialGroups []BloodGroup
	for _, bt := range bloodTypes {
		initialGroups = append(initialGroups, BloodGroup{
			Group:    bt,
			Quantity: 0, // Stock starts empty
		})
	}

	return Bank{
		Name:        name,
		PhoneNumber: phone,
		Email:       email,
		BloodGroups: initialGroups, // GORM handles creation
		LastUpdated: time.Now(),
	}
}