package models

import (
	"time"
)

type Bank struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	LastUpdated time.Time `json:"last_updated"`

	// One-to-One Relationship with Address
	// Bank has one Address
	Address Address `json:"address" gorm:"foreignKey:BankID"`

	// One-to-One Relationship with Photo
	// Bank has one Photo (Main photo)
	// We assume Photo struct has a BankID
	Photo Photo `json:"photo" gorm:"foreignKey:BankID"`

	// One-to-Many Relationship with BloodGroups (Inventory)
	// Bank has many BloodGroups (Stock)
	BloodGroups []BloodGroup `json:"blood_groups" gorm:"foreignKey:BankID"`

	// One-to-Many (or Many-to-Many) with Moderators
	// Assuming a Moderator belongs to a Bank
	Moderators []Moderator `json:"moderators" gorm:"foreignKey:BankID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBank(name, phone, email string) Bank {
	bloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	
	var initialGroups []BloodGroup
	for _, bt := range bloodTypes {
		initialGroups = append(initialGroups, BloodGroup{
			Group:    bt,
			Quantity: 0, // كيبدا بـ 0
		})
	}

	return Bank{
		Name:        name,
		PhoneNumber: phone,
		Email:       email,
		BloodGroups: initialGroups, 
		LastUpdated: time.Now(),
	}
}