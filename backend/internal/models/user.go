package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	
	// Identity Fields (Replacement for IdentityUser)
	UserName     string `gorm:"unique;not null" json:"username"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `json:"-"` // "-" means don't send in JSON response
	PhoneNumber  string `json:"phone_number"`

	// Profile Fields
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"` // "male", "female"
	LastActive  time.Time `json:"last_active"`
	
	// Blood Info
	BloodGroup string `json:"blood_group"` // e.g. "O+"
	Available  bool   `json:"available"`

	// Relationships
	Address Address `json:"address" gorm:"foreignKey:UserID"`

	Photo Photo `json:"photo" gorm:"foreignKey:UserID"`

	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
	Moderates []Moderator `json:"moderates" gorm:"foreignKey:UserID"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetAge() int {
	if u.DateOfBirth.IsZero() {
		return 0
	}
	now := time.Now()
	age := now.Year() - u.DateOfBirth.Year()
	if now.YearDay() < u.DateOfBirth.YearDay() {
		age--
	}
	return age
}