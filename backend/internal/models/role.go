package models

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique" json:"name"` // e.g. "Admin"
	
	// Users []User `gorm:"many2many:user_roles;"` 
}