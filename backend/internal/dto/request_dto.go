package dto

type CreateRequestInput struct {
	BloodType    string  `json:"bloodType" validate:"required"`
	IsUrgent     bool    `json:"isUrgent"`
	HospitalName string  `json:"hospitalName" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	// Note: Photo kan-siftoha bchkel akhor (Multipart Form)
}