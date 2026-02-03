package dto

// BloodGroupUpdateInput: 
type BloodGroupUpdateInput struct {
	BloodGroupID uint `json:"blood_group_id" validate:"required"`
	Quantity     int  `json:"quantity" validate:"required,min=0"` 
}